package inclusion

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/google/uuid"

	"github.com/educabot/alizia-inclusion-be/src/core/entities"
	"github.com/educabot/alizia-inclusion-be/src/core/providers"
)

// maxAgenticIterations caps how many tool-calling rounds a single chat turn may
// run before we force the model to answer, preventing an infinite tool loop.
// One pass of tools, two at most — keeps latency low in the classroom and bounds
// cost. Raise if usage patterns justify it.
const maxAgenticIterations = 2

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
) (*providers.ChatResponse, int, error) {
	// No tools: collapse to a single plain Chat, identical to non-agentic behavior.
	if len(tools) == 0 {
		resp, err := ai.Chat(ctx, messages)
		return resp, 0, err
	}

	var totalUsage providers.TokenUsage
	var sawUsage bool
	toolCalls := 0 // cumulative tool calls in this turn, used for tracing

	for range maxIters {
		resp, err := ai.ChatWithTools(ctx, messages, tools)
		if err != nil {
			return nil, toolCalls, err
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
			return resp, toolCalls, nil
		}
		toolCalls += len(resp.ToolCalls)

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
		return nil, toolCalls, err
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
	return final, toolCalls, nil
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
		{
			Name:        "get_student_history",
			Description: "Devuelve un resumen de las conversaciones previas sobre un alumno (de qué se venía hablando). Útil para retomar el hilo sin recontextualizar todo.",
			Parameters: map[string]any{
				"type": "object",
				"properties": map[string]any{
					"student_id": map[string]any{
						"type":        "integer",
						"description": "ID del alumno cuyo historial de conversaciones se quiere recuperar.",
					},
				},
				"required": []string{"student_id"},
			},
		},
		{
			Name:        "get_past_adaptations",
			Description: "Lista las adaptaciones previas de un alumno con su estado y resultado en aula. Útil para no repetir lo que ya se probó y construir sobre lo que funcionó.",
			Parameters: map[string]any{
				"type": "object",
				"properties": map[string]any{
					"student_id": map[string]any{
						"type":        "integer",
						"description": "ID del alumno cuyas adaptaciones previas se quieren listar.",
					},
				},
				"required": []string{"student_id"},
			},
		},
		{
			Name:        "search_content",
			Description: "Busca material pedagógico real (libros / papers / guías) sobre un tema de inclusión. Reescribí la pregunta del docente a palabras clave (temas y discapacidades, ej. 'TEA autismo autorregulación') antes de llamar. Devuelve los fragmentos más relevantes con un preview. Si vuelve vacío, no inventes: respondé con los lineamientos base aclarando que no hay material cargado.",
			Parameters: map[string]any{
				"type": "object",
				"properties": map[string]any{
					"query": map[string]any{
						"type":        "string",
						"description": "Palabras clave del tema a buscar (temas + nombres de discapacidades).",
					},
				},
				"required": []string{"query"},
			},
		},
		{
			Name:        "get_content",
			Description: "Trae el contenido pedagógico completo de un documento por su id (obtenido de search_content), con todos sus fragmentos.",
			Parameters: map[string]any{
				"type": "object",
				"properties": map[string]any{
					"content_id": map[string]any{
						"type":        "integer",
						"description": "ID del documento pedagógico a recuperar.",
					},
				},
				"required": []string{"content_id"},
			},
		},
		{
			Name:        "create_student",
			Description: "Da de alta un alumno. GUARDADO CONFIRMADO: primero llamá SIN confirmed (o confirmed=false) para armar la ficha y mostrársela al docente; recién cuando el docente confirme, volvé a llamar con confirmed=true para persistir. Nunca guardes en silencio.",
			Parameters: map[string]any{
				"type": "object",
				"properties": map[string]any{
					"name":         map[string]any{"type": "string", "description": "Nombre del alumno."},
					"classroom_id": map[string]any{"type": "integer", "description": "ID del aula."},
					"confirmed":    map[string]any{"type": "boolean", "description": "true solo cuando el docente confirmó el alta."},
				},
				"required": []string{"name", "classroom_id"},
			},
		},
		{
			Name:        "create_recurso",
			Description: "Crea un recurso (adaptación) para un alumno usando dispositivos de la valija. GUARDADO CONFIRMADO: llamá primero SIN confirmed para mostrar la ficha de frente al docente y preguntarle si la guarda; recién con confirmed=true persiste, ligada al alumno y a esta conversación. Es privada del docente.",
			Parameters: map[string]any{
				"type": "object",
				"properties": map[string]any{
					"student_id":           map[string]any{"type": "integer", "description": "ID del alumno al que se vincula."},
					"subject":              map[string]any{"type": "string", "description": "Materia o área."},
					"title":                map[string]any{"type": "string", "description": "Título corto del recurso."},
					"activity_description": map[string]any{"type": "string", "description": "Descripción de la actividad adaptada."},
					"adaptation_strategy":  map[string]any{"type": "string", "description": "Resumen de la estrategia."},
					"adaptation_type":      map[string]any{"type": "string", "description": "Tipo (actividad_adaptada, material_nuevo, estrategia_aula, situacion_emergente)."},
					"device_ids":           map[string]any{"type": "array", "items": map[string]any{"type": "integer"}, "description": "IDs de dispositivos de la valija usados."},
					"was_edited":           map[string]any{"type": "boolean", "description": "true si el docente editó la propuesta antes de guardar."},
					"confirmed":            map[string]any{"type": "boolean", "description": "true solo cuando el docente confirmó guardar."},
				},
				"required": []string{"student_id", "subject"},
			},
		},
		{
			Name:        "relate_student_recurso",
			Description: "Vincula un recurso (adaptación) existente a un alumno. GUARDADO CONFIRMADO: confirmed=true solo tras el OK del docente.",
			Parameters: map[string]any{
				"type": "object",
				"properties": map[string]any{
					"student_id": map[string]any{"type": "integer", "description": "ID del alumno."},
					"recurso_id": map[string]any{"type": "integer", "description": "ID del recurso (adaptación)."},
					"confirmed":  map[string]any{"type": "boolean", "description": "true solo cuando el docente confirmó."},
				},
				"required": []string{"student_id", "recurso_id"},
			},
		},
	}
}

// inclusionDispatcher executes inclusionTools against the domain providers.
// summaries and adaptations are optional: missing providers return a handled
// error instead of panicking.
type inclusionDispatcher struct {
	students    providers.StudentProvider
	devices     providers.DeviceProvider
	summaries   providers.ConversationSummaryProvider
	adaptations providers.AdaptationProvider
	content     providers.PedagogicalContentProvider
	// Turn identity for action tools: who is saving and which conversation
	// originated the resource.
	userID         int64
	conversationID int64
}

// defaultContentSearchLimit caps how many chunks the RAG layer returns per search.
const defaultContentSearchLimit = 5

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

	case "get_student_history":
		if d.summaries == nil {
			return "", fmt.Errorf("get_student_history no disponible")
		}
		var args struct {
			StudentID int64 `json:"student_id"`
		}
		if err := json.Unmarshal([]byte(call.Arguments), &args); err != nil {
			return "", fmt.Errorf("invalid arguments for get_student_history: %w", err)
		}
		summaries, err := d.summaries.RecentByStudent(ctx, orgID, args.StudentID, maxPriorSummaries)
		if err != nil {
			return "", err
		}
		type summaryLite struct {
			ConversationID int64    `json:"conversation_id"`
			Summary        string   `json:"summary"`
			TopicKeywords  []string `json:"topic_keywords,omitempty"`
		}
		out := make([]summaryLite, len(summaries))
		for i := range summaries {
			out[i] = summaryLite{
				ConversationID: summaries[i].ConversationID,
				Summary:        summaries[i].Summary,
				TopicKeywords:  summaries[i].TopicKeywords,
			}
		}
		return marshalToolResult(map[string]any{"history": out})

	case "get_past_adaptations":
		if d.adaptations == nil {
			return "", fmt.Errorf("get_past_adaptations no disponible")
		}
		var args struct {
			StudentID int64 `json:"student_id"`
		}
		if err := json.Unmarshal([]byte(call.Arguments), &args); err != nil {
			return "", fmt.Errorf("invalid arguments for get_past_adaptations: %w", err)
		}
		adaptations, err := d.adaptations.List(ctx, orgID, providers.AdaptationFilter{StudentID: &args.StudentID})
		if err != nil {
			return "", err
		}
		type adaptationLite struct {
			ID      int64  `json:"id"`
			Subject string `json:"subject"`
			Status  string `json:"status"`
			Outcome string `json:"outcome,omitempty"`
		}
		out := make([]adaptationLite, len(adaptations))
		for i := range adaptations {
			lite := adaptationLite{ID: adaptations[i].ID, Subject: adaptations[i].Subject, Status: adaptations[i].Status}
			if adaptations[i].Outcome != nil {
				lite.Outcome = *adaptations[i].Outcome
			}
			out[i] = lite
		}
		return marshalToolResult(map[string]any{"adaptations": out})

	case "search_content":
		if d.content == nil {
			return "", fmt.Errorf("search_content no disponible")
		}
		var args struct {
			Query string `json:"query"`
		}
		if err := json.Unmarshal([]byte(call.Arguments), &args); err != nil {
			return "", fmt.Errorf("invalid arguments for search_content: %w", err)
		}
		results, err := d.content.SearchChunks(ctx, orgID, args.Query, defaultContentSearchLimit)
		if err != nil {
			return "", err
		}
		// Empty results are returned explicitly so the LLM falls back to base
		// guidelines instead of hallucinating content.
		return marshalToolResult(map[string]any{"results": results})

	case "get_content":
		if d.content == nil {
			return "", fmt.Errorf("get_content no disponible")
		}
		var args struct {
			ContentID int64 `json:"content_id"`
		}
		if err := json.Unmarshal([]byte(call.Arguments), &args); err != nil {
			return "", fmt.Errorf("invalid arguments for get_content: %w", err)
		}
		content, err := d.content.GetContent(ctx, orgID, args.ContentID)
		if err != nil {
			return "", err
		}
		return marshalToolResult(content)

	case "create_student":
		if d.students == nil {
			return "", fmt.Errorf("create_student no disponible")
		}
		var args struct {
			Name        string `json:"name"`
			ClassroomID int64  `json:"classroom_id"`
			Confirmed   bool   `json:"confirmed"`
		}
		if err := json.Unmarshal([]byte(call.Arguments), &args); err != nil {
			return "", fmt.Errorf("invalid arguments for create_student: %w", err)
		}
		student := &entities.Student{OrganizationID: orgID, ClassroomID: args.ClassroomID, Name: args.Name}
		if !args.Confirmed {
			// Dry-run: return the draft record so the LLM can display it and request confirmation.
			return marshalToolResult(map[string]any{"pending_confirmation": true, "student": student})
		}
		if err := d.students.Create(ctx, student); err != nil {
			return "", err
		}
		return marshalToolResult(map[string]any{"saved": true, "student": student})

	case "create_recurso":
		if d.adaptations == nil {
			return "", fmt.Errorf("create_recurso no disponible")
		}
		var args struct {
			StudentID           int64   `json:"student_id"`
			Subject             string  `json:"subject"`
			Title               string  `json:"title"`
			ActivityDescription *string `json:"activity_description"`
			AdaptationStrategy  *string `json:"adaptation_strategy"`
			AdaptationType      string  `json:"adaptation_type"`
			DeviceIDs           []int64 `json:"device_ids"`
			WasEdited           bool    `json:"was_edited"`
			Confirmed           bool    `json:"confirmed"`
		}
		if err := json.Unmarshal([]byte(call.Arguments), &args); err != nil {
			return "", fmt.Errorf("invalid arguments for create_recurso: %w", err)
		}
		adaptationType := args.AdaptationType
		if adaptationType == "" {
			adaptationType = defaultAdaptationType
		}
		recurso := &entities.Adaptation{
			OrganizationID:      orgID,
			StudentID:           args.StudentID,
			TeacherID:           d.userID,
			Title:               args.Title,
			Subject:             args.Subject,
			ActivityDescription: args.ActivityDescription,
			AdaptationStrategy:  args.AdaptationStrategy,
			AdaptationType:      adaptationType,
			Status:              "en_curso",
			WasEdited:           args.WasEdited,
		}
		if d.conversationID > 0 {
			recurso.SourceConversationID = &d.conversationID
		}
		if !args.Confirmed {
			// Dry-run: return the draft record; nothing is persisted until the teacher confirms.
			return marshalToolResult(map[string]any{"pending_confirmation": true, "recurso": recurso, "device_ids": args.DeviceIDs})
		}
		if err := d.adaptations.Create(ctx, recurso); err != nil {
			return "", err
		}
		if len(args.DeviceIDs) > 0 {
			if err := d.adaptations.SetDevices(ctx, recurso.ID, args.DeviceIDs); err != nil {
				return "", err
			}
		}
		return marshalToolResult(map[string]any{"saved": true, "recurso": recurso})

	case "relate_student_recurso":
		if d.adaptations == nil {
			return "", fmt.Errorf("relate_student_recurso no disponible")
		}
		var args struct {
			StudentID int64 `json:"student_id"`
			RecursoID int64 `json:"recurso_id"`
			Confirmed bool  `json:"confirmed"`
		}
		if err := json.Unmarshal([]byte(call.Arguments), &args); err != nil {
			return "", fmt.Errorf("invalid arguments for relate_student_recurso: %w", err)
		}
		recurso, err := d.adaptations.Get(ctx, orgID, args.RecursoID)
		if err != nil {
			return "", err
		}
		recurso.StudentID = args.StudentID
		if !args.Confirmed {
			return marshalToolResult(map[string]any{"pending_confirmation": true, "recurso": recurso})
		}
		if err := d.adaptations.Update(ctx, recurso); err != nil {
			return "", err
		}
		return marshalToolResult(map[string]any{"saved": true, "recurso": recurso})

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
