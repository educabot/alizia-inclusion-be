package inclusion

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/google/uuid"

	"github.com/educabot/alizia-inclusion-be/src/core/providers"
)

// maxAgenticIterations caps how many tool-calling rounds a single chat turn may
// run before we force the model to answer, preventing an infinite tool loop.
const maxAgenticIterations = 4

// toolDispatcher executes a tool the model asked for and returns its result as a
// JSON string that is fed back to the model.
type toolDispatcher interface {
	Dispatch(ctx context.Context, orgID uuid.UUID, call providers.ToolCall) (string, error)
}

// runAgenticChat drives a tool-calling loop. It calls ChatWithTools; if the model
// requests tools, it executes them via dispatcher, appends the results, and loops
// again until the model answers with no tool calls or maxIters is reached.
//
// Token usage is accumulated across every round so callers can record the full
// cost of the turn. When tools is empty the loop collapses to a single call,
// behaving identically to a plain Chat.
func runAgenticChat(
	ctx context.Context,
	ai providers.AIClient,
	messages []providers.ChatMessage,
	tools []providers.ToolDefinition,
	dispatcher toolDispatcher,
	orgID uuid.UUID,
	maxIters int,
) (*providers.ChatResponse, error) {
	// No tools: collapse to a single plain Chat, identical to non-agentic behavior.
	if len(tools) == 0 {
		return ai.Chat(ctx, messages)
	}

	var totalUsage providers.TokenUsage
	var sawUsage bool

	for range maxIters {
		resp, err := ai.ChatWithTools(ctx, messages, tools)
		if err != nil {
			return nil, err
		}
		if resp.Usage != nil {
			sawUsage = true
			totalUsage.PromptTokens += resp.Usage.PromptTokens
			totalUsage.CompletionTokens += resp.Usage.CompletionTokens
			totalUsage.TotalTokens += resp.Usage.TotalTokens
		}

		if len(resp.ToolCalls) == 0 {
			if sawUsage {
				resp.Usage = &totalUsage
			}
			return resp, nil
		}

		// Echo the assistant turn (with its tool calls) so the model keeps context.
		messages = append(messages, providers.ChatMessage{
			Role:      "assistant",
			Content:   resp.Content,
			ToolCalls: resp.ToolCalls,
		})

		// Execute each requested tool and append its result as a tool message.
		for _, call := range resp.ToolCalls {
			result, derr := dispatcher.Dispatch(ctx, orgID, call)
			if derr != nil {
				result = fmt.Sprintf(`{"error":%q}`, derr.Error())
			}
			messages = append(messages, providers.ChatMessage{
				Role:       "tool",
				Content:    result,
				ToolCallID: call.ID,
			})
		}
	}

	// Iteration budget exhausted: ask once more without tools to force an answer.
	final, err := ai.Chat(ctx, messages)
	if err != nil {
		return nil, err
	}
	if final.Usage != nil {
		sawUsage = true
		totalUsage.PromptTokens += final.Usage.PromptTokens
		totalUsage.CompletionTokens += final.Usage.CompletionTokens
		totalUsage.TotalTokens += final.Usage.TotalTokens
	}
	if sawUsage {
		final.Usage = &totalUsage
	}
	return final, nil
}

// inclusionTools are the domain tools Alizia can call to ground its answers in
// real classroom data instead of relying solely on the prompt context.
func inclusionTools() []providers.ToolDefinition {
	return []providers.ToolDefinition{
		{
			Name:        "list_classroom_students",
			Description: "Lista los alumnos de un aula con su id y nombre. Útil para identificar a quién se refiere el docente.",
			Parameters: map[string]any{
				"type": "object",
				"properties": map[string]any{
					"classroom_id": map[string]any{
						"type":        "integer",
						"description": "ID del aula cuyos alumnos se quieren listar.",
					},
				},
				"required": []string{"classroom_id"},
			},
		},
		{
			Name:        "get_student",
			Description: "Devuelve los datos y el perfil de un alumno por su id, incluyendo necesidades de apoyo.",
			Parameters: map[string]any{
				"type": "object",
				"properties": map[string]any{
					"student_id": map[string]any{
						"type":        "integer",
						"description": "ID del alumno.",
					},
				},
				"required": []string{"student_id"},
			},
		},
		{
			Name:        "list_devices",
			Description: "Lista los dispositivos de la valija adaptativa disponibles, con su id, nombre y para qué sirven.",
			Parameters: map[string]any{
				"type":       "object",
				"properties": map[string]any{},
			},
		},
	}
}

// inclusionDispatcher executes inclusionTools against the domain providers.
type inclusionDispatcher struct {
	students providers.StudentProvider
	devices  providers.DeviceProvider
}

func (d inclusionDispatcher) Dispatch(ctx context.Context, orgID uuid.UUID, call providers.ToolCall) (string, error) {
	switch call.Name {
	case "list_classroom_students":
		var args struct {
			ClassroomID int64 `json:"classroom_id"`
		}
		if err := json.Unmarshal([]byte(call.Arguments), &args); err != nil {
			return "", fmt.Errorf("invalid arguments for list_classroom_students: %w", err)
		}
		students, err := d.students.ListByClassroom(ctx, orgID, args.ClassroomID)
		if err != nil {
			return "", err
		}
		type studentLite struct {
			ID   int64  `json:"id"`
			Name string `json:"name"`
		}
		out := make([]studentLite, len(students))
		for i := range students {
			out[i] = studentLite{ID: students[i].ID, Name: students[i].Name}
		}
		return marshalToolResult(map[string]any{"students": out})

	case "get_student":
		var args struct {
			StudentID int64 `json:"student_id"`
		}
		if err := json.Unmarshal([]byte(call.Arguments), &args); err != nil {
			return "", fmt.Errorf("invalid arguments for get_student: %w", err)
		}
		student, err := d.students.GetStudent(ctx, orgID, args.StudentID)
		if err != nil {
			return "", err
		}
		return marshalToolResult(student)

	case "list_devices":
		devices, err := d.devices.ListDevices(ctx, orgID, nil)
		if err != nil {
			return "", err
		}
		type deviceLite struct {
			ID         int64  `json:"id"`
			Name       string `json:"name"`
			UsefulWhen string `json:"useful_when,omitempty"`
		}
		out := make([]deviceLite, len(devices))
		for i := range devices {
			lite := deviceLite{ID: devices[i].ID, Name: devices[i].Name}
			if devices[i].UsefulWhen != nil {
				lite.UsefulWhen = *devices[i].UsefulWhen
			}
			out[i] = lite
		}
		return marshalToolResult(map[string]any{"devices": out})

	default:
		return "", fmt.Errorf("unknown tool: %s", call.Name)
	}
}

func marshalToolResult(v any) (string, error) {
	b, err := json.Marshal(v)
	if err != nil {
		return "", fmt.Errorf("marshal tool result: %w", err)
	}
	return string(b), nil
}
