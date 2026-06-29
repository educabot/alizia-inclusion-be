package inclusion

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"strings"

	"github.com/google/uuid"

	"github.com/educabot/alizia-inclusion-be/src/core/providers"
	"github.com/educabot/alizia-inclusion-be/src/observability"
)

// maxAgenticIterations caps how many tool-calling rounds a single chat turn may
// run before we force the model to answer, preventing an infinite tool loop.
const maxAgenticIterations = 4

// toolDispatcher executes a tool the model asked for and returns its result as a
// JSON string that is fed back to the model.
type toolDispatcher interface {
	Dispatch(ctx context.Context, orgID uuid.UUID, call providers.ToolCall) (string, error)
}

// toolInvocation registra una tool ejecutada en el turno: para trazabilidad
// (chat.sources_used) y para poblar referenced_content con el material citado.
type toolInvocation struct {
	Name        string
	StudentID   *int64
	Query       string
	Hits        int
	ContentRefs []ContentRef
	Failed      bool
}

// runAgenticChat drives a tool-calling loop. It calls ChatWithTools; if the model
// requests tools, it executes them via dispatcher, appends the results, and loops
// again until the model answers with no tool calls or maxIters is reached.
//
// Token usage is accumulated across every round so callers can record the full
// cost of the turn. When tools is empty the loop collapses to a single call,
// behaving identically to a plain Chat. Devuelve además el trace de tools
// ejecutadas (vacío en el camino sin tools) para que el caller derive las fuentes
// usadas y los contenidos citados.
func runAgenticChat(
	ctx context.Context,
	ai providers.AIClient,
	messages []providers.ChatMessage,
	tools []providers.ToolDefinition,
	dispatcher toolDispatcher,
	orgID uuid.UUID,
	maxIters int,
) (*providers.ChatResponse, []toolInvocation, error) {
	// No tools: collapse to a single plain Chat, identical to non-agentic behavior.
	if len(tools) == 0 {
		resp, err := ai.Chat(ctx, messages)
		return resp, nil, err
	}

	var totalUsage providers.TokenUsage
	var sawUsage bool
	var trace []toolInvocation

	for range maxIters {
		resp, err := ai.ChatWithTools(ctx, messages, tools)
		if err != nil {
			return nil, trace, err
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
			return resp, trace, nil
		}

		slog.InfoContext(ctx, "chat.agentic_iteration", "tool_calls", len(resp.ToolCalls))

		// Echo the assistant turn (with its tool calls) so the model keeps context.
		messages = append(messages, providers.ChatMessage{
			Role:      "assistant",
			Content:   resp.Content,
			ToolCalls: resp.ToolCalls,
		})

		// Execute each requested tool and append its result as a tool message.
		for _, call := range resp.ToolCalls {
			slog.InfoContext(ctx, "chat.tool_call", "tool", call.Name, observability.Text("args", call.Arguments))
			result, derr := dispatcher.Dispatch(ctx, orgID, call)
			inv := toolInvocation{
				Name:      call.Name,
				StudentID: extractToolStudentID(call.Arguments),
				Query:     extractToolQuery(call.Arguments),
			}
			if derr != nil {
				inv.Failed = true
				result = fmt.Sprintf(`{"error":%q}`, derr.Error())
				slog.WarnContext(ctx, "chat.tool_error", "tool", call.Name, "error", derr.Error())
			} else {
				inv.Hits = countResults(result)
				inv.ContentRefs = extractContentRefs(call.Name, result)
				slog.InfoContext(ctx, "chat.tool_result", "tool", call.Name, "result_len", len(result), observability.Text("result", result))
			}
			trace = append(trace, inv)
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
		return nil, trace, err
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
	return final, trace, nil
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
			Description: "Busca contenido pedagógico por texto en la base de materiales del corpus (búsqueda full-text clásica). Usá search_content_hibrido cuando quieras resultados por similitud semántica; usá search_content para búsquedas exactas por término.",
			Parameters: map[string]any{
				"type": "object",
				"properties": map[string]any{
					"query": map[string]any{
						"type":        "string",
						"description": "Términos de búsqueda a buscar en el contenido pedagógico.",
					},
				},
				"required": []string{"query"},
			},
		},
		{
			Name:        "search_content_hibrido",
			Description: "Búsqueda semántica híbrida (vector + texto + conceptos) en el corpus de material pedagógico. Pasá la pregunta del docente COMPLETA en semantic_question (se usa para el embedding) y, opcionalmente, palabras clave en terms para reforzar. Devuelve los fragmentos más relevantes con score, fuente y un extracto. Si vuelve vacío, no inventes.",
			Parameters: map[string]any{
				"type": "object",
				"properties": map[string]any{
					"semantic_question": map[string]any{
						"type":        "string",
						"description": "La pregunta del docente completa, en lenguaje natural.",
					},
					"terms": map[string]any{
						"type":        "array",
						"items":       map[string]any{"type": "string"},
						"description": "Palabras o locuciones clave para reforzar (temas, discapacidades, nombres de guías).",
					},
					"resource_id": map[string]any{
						"type":        "integer",
						"description": "Opcional: acota la búsqueda a un documento puntual por su id.",
					},
					"limit": map[string]any{
						"type":        "integer",
						"description": "Máximo de fragmentos a devolver (default 5).",
					},
				},
				"required": []string{"semantic_question"},
			},
		},
		{
			Name:        "get_content",
			Description: "Trae el contenido pedagógico completo de un documento por su id (obtenido de search_content_hibrido), con todos sus fragmentos.",
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
	}
}

// inclusionDispatcher executes inclusionTools against the domain providers.
// summaries y adaptations son opcionales: si faltan, sus tools devuelven un
// error manejable en vez de panicar.
type inclusionDispatcher struct {
	students    providers.StudentProvider
	devices     providers.DeviceProvider
	summaries   providers.ConversationSummaryProvider
	adaptations providers.AdaptationProvider
	content     providers.PedagogicalContentProvider
	embedder    providers.Embedder
	rag         providers.RAGSearchProvider
}

// defaultContentSearchLimit acota cuántos chunks devuelve el RAG por búsqueda.
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
		adaptations, err := d.adaptations.List(ctx, orgID, &args.StudentID, nil)
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
		// Sin coincidencias: devolvemos lista vacía explícita para que la LLM
		// caiga a los lineamientos base sin inventar.
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

	case "search_content_hibrido":
		if d.rag == nil || d.embedder == nil {
			return "", fmt.Errorf("search_content_hibrido no disponible")
		}
		var args struct {
			SemanticQuestion string   `json:"semantic_question"`
			Terms            []string `json:"terms"`
			ResourceID       *int64   `json:"resource_id"`
			Limit            int      `json:"limit"`
		}
		if err := json.Unmarshal([]byte(call.Arguments), &args); err != nil {
			return "", fmt.Errorf("invalid arguments for search_content_hibrido: %w", err)
		}
		limit := args.Limit
		if limit <= 0 {
			limit = defaultContentSearchLimit
		}
		embedding, err := d.embedder.EmbedQuery(ctx, args.SemanticQuestion)
		if err != nil {
			return "", err
		}
		slog.InfoContext(ctx, "rag.embed_ok", "question", args.SemanticQuestion, "dims", len(embedding), "terms", args.Terms)
		hits, err := d.rag.HybridSearch(ctx, providers.HybridSearchSpec{
			ResourceID:       args.ResourceID,
			SemanticQuestion: args.SemanticQuestion,
			Terms:            args.Terms,
			Limit:            limit,
		}, embedding)
		if err != nil {
			return "", err
		}
		// Vista recortada para la LLM: sin el content completo, solo un extracto.
		type chunkLite struct {
			ResourceID int64    `json:"resource_id"`
			Title      string   `json:"title"`
			Score      float64  `json:"score"`
			Pages      string   `json:"pages,omitempty"`
			Summary    string   `json:"summary,omitempty"`
			Concepts   []string `json:"concepts,omitempty"`
			Snippet    string   `json:"snippet"`
		}
		out := make([]chunkLite, len(hits))
		for i := range hits {
			out[i] = chunkLite{
				ResourceID: hits[i].ResourceID,
				Title:      hits[i].Title,
				Score:      hits[i].Score,
				Pages:      fmt.Sprintf("%d-%d", hits[i].PageStart, hits[i].PageEnd),
				Summary:    hits[i].Summary,
				Concepts:   hits[i].Concepts,
				Snippet:    snippet(hits[i].Content, 650),
			}
		}
		return marshalToolResult(map[string]any{"results": out})

	default:
		return "", fmt.Errorf("unknown tool: %s", call.Name)
	}
}

// snippet normaliza espacios y recorta a maxRunes para no inflar el contexto de la LLM.
func snippet(text string, maxRunes int) string {
	clean := strings.Join(strings.Fields(text), " ")
	runes := []rune(clean)
	if len(runes) <= maxRunes {
		return clean
	}
	return string(runes[:maxRunes]) + "..."
}

func marshalToolResult(v any) (string, error) {
	b, err := json.Marshal(v)
	if err != nil {
		return "", fmt.Errorf("marshal tool result: %w", err)
	}
	return string(b), nil
}
