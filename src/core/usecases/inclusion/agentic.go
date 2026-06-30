package inclusion

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"strings"

	"github.com/google/uuid"

	"github.com/educabot/alizia-inclusion-be/src/core/entities"
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
	requireSearchBeforeProposal bool,
) (*providers.ChatResponse, []toolInvocation, error) {
	// No tools: collapse to a single plain Chat, identical to non-agentic behavior.
	if len(tools) == 0 {
		resp, err := ai.Chat(ctx, messages)
		return resp, nil, err
	}

	var totalUsage providers.TokenUsage
	var sawUsage bool
	var trace []toolInvocation
	// forcedSearch evita forzar la búsqueda más de una vez (no loopear si el modelo
	// se empecina en proponer sin buscar): la red de seguridad actúa una sola vez.
	var forcedSearch bool

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
			// Red de seguridad RAG: si el turno final emite una propuesta (paso a paso o
			// recurso) pero en toda la conversación nunca se buscó en la bibliografía,
			// forzamos una búsqueda y le pedimos reformular fundamentando. Mecanismo
			// general (no regla por caso), una sola vez, dentro del presupuesto de iteraciones.
			if requireSearchBeforeProposal && !forcedSearch && emitsProposal(resp.Content) && !traceHasTool(trace, "search_content_hibrido") {
				forcedSearch = true
				slog.WarnContext(ctx, "chat.rag_safety_net_triggered", "reason", "proposal_without_search")
				messages = append(messages,
					providers.ChatMessage{Role: "assistant", Content: resp.Content},
					providers.ChatMessage{Role: "user", Content: "Antes de darme esta propuesta tenés que fundamentarla: llamá a search_content_hibrido con la barrera observable y la edad, y recién entonces reformulá integrando lo que encuentres. No propongas un paso a paso ni guardes el recurso sin haber buscado en la bibliografía."},
				)
				continue
			}
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

// emitsProposal indica si el contenido visible es una propuesta accionable: un paso a
// paso ([STEPS]) o un recurso a guardar ([ADAPTATION_JSON]). Es la señal de que el
// modelo "cerró" con una recomendación y, por ende, debería haberse fundamentado.
func emitsProposal(content string) bool {
	return strings.Contains(content, "[STEPS]") || strings.Contains(content, "[ADAPTATION_JSON")
}

// traceHasTool indica si alguna tool con ese nombre se ejecutó (sin error) en el turno.
func traceHasTool(trace []toolInvocation, name string) bool {
	for i := range trace {
		if trace[i].Name == name && !trace[i].Failed {
			return true
		}
	}
	return false
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
		{
			Name:        "find_student_by_name",
			Description: "Busca alumnos por nombre (aproximado) en toda la organización, más allá de la lista de alumnos que conocés. Úsala para reconocer a un alumno ANTES de ofrecer crearlo. Devuelve los que coinciden con su id, aula y dificultades.",
			Parameters: map[string]any{
				"type": "object",
				"properties": map[string]any{
					"name": map[string]any{
						"type":        "string",
						"description": "Nombre o parte del nombre del alumno a buscar.",
					},
				},
				"required": []string{"name"},
			},
		},
		{
			Name:        "list_classrooms",
			Description: "Lista las aulas de la organización con su id, nombre, grado y sección. Útil para ubicar el aula de un alumno antes de crearlo.",
			Parameters: map[string]any{
				"type":       "object",
				"properties": map[string]any{},
			},
		},
		{
			Name:        "create_classroom",
			Description: "Crea un aula nueva a partir de cómo la nombra el docente ('3ro A', 'tercero B'). Úsala solo si el aula no existe (revisá antes con list_classrooms). Devuelve el aula con su id, que usás como classroom_id en create_student.",
			Parameters: map[string]any{
				"type": "object",
				"properties": map[string]any{
					"grade": map[string]any{
						"type":        "string",
						"description": "El aula tal como la nombra el docente, ej. '3ro A' o 'tercero B'.",
					},
				},
				"required": []string{"grade"},
			},
		},
		{
			Name:        "create_student",
			Description: "Da de alta un alumno NUEVO. El aula (classroom_id) es OPCIONAL: si el docente todavía no la dijo, creá igual al alumno sin aula y pedísela después para asentarla (no es un bloqueante). Usala cuando el docente nombra a un alumno concreto y le estás armando un recurso, para dejarlo asociado; no la uses para alumnos que ya reconociste con find_student_by_name. Es idempotente por nombre: si el alumno ya existe lo devuelve (y si ahora le pasás classroom_id y no tenía aula, se la fija). Devuelve el alumno con su id, que usás como [STUDENT_ID:X] y para enlazar la adaptación. Para fijar el aula más tarde, rellamá esta tool con el mismo name + el classroom_id (resuelto con list_classrooms / create_classroom).",
			Parameters: map[string]any{
				"type": "object",
				"properties": map[string]any{
					"name": map[string]any{
						"type":        "string",
						"description": "Nombre completo del alumno, tal como lo dio el docente.",
					},
					"classroom_id": map[string]any{
						"type":        "integer",
						"description": "OPCIONAL: ID del aula (de list_classrooms o create_classroom). Si todavía no la sabés, omitilo: el alumno se crea sin aula y la completás después.",
					},
					"preferred_name": map[string]any{
						"type":        "string",
						"description": "Opcional: cómo se le llama habitualmente.",
					},
					"difficulties": map[string]any{
						"type":        "array",
						"items":       map[string]any{"type": "string"},
						"description": "Opcional: necesidades para el aprendizaje observables en el aula (no diagnósticos), ej. 'le cuesta sostener la atención'.",
					},
					"is_transitory": map[string]any{
						"type":        "boolean",
						"description": "Opcional: true si la condición es transitoria.",
					},
					"free_description": map[string]any{
						"type":        "string",
						"description": "Opcional: descripción libre del contexto del alumno.",
					},
				},
				"required": []string{"name"},
			},
		},
	}
}

// inclusionDispatcher executes inclusionTools against the domain providers.
// summaries y adaptations son opcionales: si faltan, sus tools devuelven un
// error manejable en vez de panicar.
type inclusionDispatcher struct {
	students    providers.StudentProvider
	profiles    providers.StudentProfileProvider
	classrooms  providers.ClassroomProvider
	devices     providers.DeviceProvider
	summaries   providers.ConversationSummaryProvider
	adaptations providers.AdaptationProvider
	content     providers.PedagogicalContentProvider
	embedder    providers.Embedder
	rag         providers.RAGSearchProvider
	// userID es el docente del turno: orgID y userID se inyectan del request, nunca
	// los pone la LLM. Se usa para que las adaptaciones que ve sean solo las suyas.
	userID int64
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
		// Solo las adaptaciones del propio docente (recursos privados).
		filter := providers.AdaptationFilter{OrgID: orgID, StudentID: &args.StudentID}
		if d.userID > 0 {
			filter.TeacherID = &d.userID
		}
		adaptations, err := d.adaptations.List(ctx, filter)
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

	case "create_student":
		var args struct {
			Name            string   `json:"name"`
			ClassroomID     int64    `json:"classroom_id"`
			PreferredName   string   `json:"preferred_name"`
			IsTransitory    bool     `json:"is_transitory"`
			Difficulties    []string `json:"difficulties"`
			FreeDescription string   `json:"free_description"`
		}
		if err := json.Unmarshal([]byte(call.Arguments), &args); err != nil {
			return "", fmt.Errorf("invalid arguments for create_student: %w", err)
		}
		name := strings.TrimSpace(args.Name)
		if name == "" {
			return "", fmt.Errorf("create_student: el nombre es obligatorio")
		}
		// El aula es OPCIONAL: el alumno se puede crear sin aula (no bloqueante). Si llega,
		// la usamos; si no, queda nil (columna NULL) y se completa después.
		var classroomID *int64
		if args.ClassroomID > 0 {
			cid := args.ClassroomID
			classroomID = &cid
		}
		// Idempotencia preservando el aula (que sigue siendo la categoría normal):
		//  - Con aula: buscamos en ESA aula (findStudentInClassroom, como siempre). Si no está,
		//    vemos si hay un alumno "sin aula" con ese nombre (creado antes) para ASENTARLE el
		//    aula (backfill) en vez de duplicarlo.
		//  - Sin aula: deduplicamos entre los "sin aula" por nombre.
		var existing *entities.Student
		var err error
		if classroomID != nil {
			existing, err = d.findStudentInClassroom(ctx, orgID, *classroomID, name)
			if err != nil {
				return "", err
			}
			if existing == nil {
				unassigned, uerr := d.findUnassignedStudentByName(ctx, orgID, name)
				if uerr != nil {
					return "", uerr
				}
				if unassigned != nil {
					unassigned.ClassroomID = classroomID
					if err := d.students.Update(ctx, unassigned); err != nil {
						return "", err
					}
					return marshalToolResult(unassigned)
				}
			}
		} else {
			existing, err = d.findUnassignedStudentByName(ctx, orgID, name)
			if err != nil {
				return "", err
			}
		}
		if existing != nil {
			return marshalToolResult(existing)
		}
		student := &entities.Student{
			OrganizationID: orgID,
			ClassroomID:    classroomID,
			Name:           name,
		}
		if pn := strings.TrimSpace(args.PreferredName); pn != "" {
			student.PreferredName = &pn
		}
		if err := d.students.Create(ctx, student); err != nil {
			return "", err
		}
		// Perfil opcional: solo si el docente aportó barrera observable o contexto.
		if d.profiles != nil && (len(args.Difficulties) > 0 || strings.TrimSpace(args.FreeDescription) != "") {
			profile := &entities.StudentProfile{
				StudentID:    student.ID,
				IsTransitory: args.IsTransitory,
				Difficulties: args.Difficulties,
			}
			if fd := strings.TrimSpace(args.FreeDescription); fd != "" {
				profile.FreeDescription = &fd
			}
			if err := d.profiles.Upsert(ctx, profile); err != nil {
				return "", err
			}
			student.Profile = profile
		}
		return marshalToolResult(student)

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

	case "find_student_by_name":
		var args struct {
			Name string `json:"name"`
		}
		if err := json.Unmarshal([]byte(call.Arguments), &args); err != nil {
			return "", fmt.Errorf("invalid arguments for find_student_by_name: %w", err)
		}
		query := normalizeName(args.Name)
		if query == "" {
			return "", fmt.Errorf("find_student_by_name: el nombre es obligatorio")
		}
		students, err := d.students.List(ctx, orgID)
		if err != nil {
			return "", err
		}
		type studentMatch struct {
			ID           int64    `json:"id"`
			Name         string   `json:"name"`
			ClassroomID  int64    `json:"classroom_id"`
			GradeLevel   string   `json:"grade_level,omitempty"`
			Difficulties []string `json:"difficulties,omitempty"`
		}
		var matches []studentMatch
		for i := range students {
			s := &students[i]
			norm := normalizeName(s.Name)
			if !strings.Contains(norm, query) {
				continue
			}
			m := studentMatch{ID: s.ID, Name: s.Name}
			if s.ClassroomID != nil {
				m.ClassroomID = *s.ClassroomID
			}
			if s.GradeLevel != nil {
				m.GradeLevel = *s.GradeLevel
			}
			if s.Profile != nil {
				m.Difficulties = s.Profile.Difficulties
			}
			matches = append(matches, m)
		}
		return marshalToolResult(map[string]any{"students": matches})

	case "list_classrooms":
		if d.classrooms == nil {
			return "", fmt.Errorf("list_classrooms no disponible")
		}
		classrooms, err := d.classrooms.List(ctx, orgID)
		if err != nil {
			return "", err
		}
		type classroomLite struct {
			ID      int64  `json:"id"`
			Name    string `json:"name"`
			Grade   string `json:"grade,omitempty"`
			Section string `json:"section,omitempty"`
		}
		out := make([]classroomLite, len(classrooms))
		for i := range classrooms {
			lite := classroomLite{ID: classrooms[i].ID, Name: classrooms[i].Name}
			if classrooms[i].Grade != nil {
				lite.Grade = *classrooms[i].Grade
			}
			if classrooms[i].Section != nil {
				lite.Section = *classrooms[i].Section
			}
			out[i] = lite
		}
		return marshalToolResult(map[string]any{"classrooms": out})

	case "create_classroom":
		if d.classrooms == nil {
			return "", fmt.Errorf("create_classroom no disponible")
		}
		var args struct {
			Grade string `json:"grade"`
		}
		if err := json.Unmarshal([]byte(call.Arguments), &args); err != nil {
			return "", fmt.Errorf("invalid arguments for create_classroom: %w", err)
		}
		name, grade, section := normalizeGrade(args.Grade)
		if name == "" {
			return "", fmt.Errorf("create_classroom: no entendí el grado (ej. '3ro A' o 'tercero B')")
		}
		// Idempotente: si ya existe un aula con ese nombre, la devolvemos.
		if existing, err := d.findClassroomByName(ctx, orgID, name); err != nil {
			return "", err
		} else if existing != nil {
			return marshalToolResult(existing)
		}
		classroom := &entities.Classroom{OrganizationID: orgID, Name: name}
		if grade != "" {
			classroom.Grade = &grade
		}
		if section != "" {
			classroom.Section = &section
		}
		if err := d.classrooms.Create(ctx, classroom); err != nil {
			return "", err
		}
		return marshalToolResult(classroom)

	default:
		return "", fmt.Errorf("unknown tool: %s", call.Name)
	}
}

// findStudentInClassroom busca un alumno por nombre normalizado dentro de un aula.
// Devuelve nil si no existe (para que create_student lo dé de alta).
func (d inclusionDispatcher) findStudentInClassroom(ctx context.Context, orgID uuid.UUID, classroomID int64, name string) (*entities.Student, error) {
	students, err := d.students.ListByClassroom(ctx, orgID, classroomID)
	if err != nil {
		return nil, err
	}
	target := normalizeName(name)
	for i := range students {
		if normalizeName(students[i].Name) == target {
			return &students[i], nil
		}
	}
	return nil, nil
}

// findUnassignedStudentByName busca un alumno SIN aula (ClassroomID nil) por nombre
// normalizado en la organización. Lo usa create_student para: (a) deduplicar cuando se crea
// sin aula, y (b) detectar un alumno creado antes sin aula para asentarle el aula después.
// Devuelve nil si no existe.
func (d inclusionDispatcher) findUnassignedStudentByName(ctx context.Context, orgID uuid.UUID, name string) (*entities.Student, error) {
	students, err := d.students.List(ctx, orgID)
	if err != nil {
		return nil, err
	}
	target := normalizeName(name)
	for i := range students {
		if students[i].ClassroomID == nil && normalizeName(students[i].Name) == target {
			return &students[i], nil
		}
	}
	return nil, nil
}

// findClassroomByName busca un aula por nombre normalizado. Devuelve nil si no existe.
func (d inclusionDispatcher) findClassroomByName(ctx context.Context, orgID uuid.UUID, name string) (*entities.Classroom, error) {
	classrooms, err := d.classrooms.List(ctx, orgID)
	if err != nil {
		return nil, err
	}
	target := normalizeName(name)
	for i := range classrooms {
		if normalizeName(classrooms[i].Name) == target {
			return &classrooms[i], nil
		}
	}
	return nil, nil
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
