package inclusion

import (
	"context"
	"fmt"
	"log/slog"
	"strings"

	"github.com/google/uuid"

	"github.com/educabot/alizia-inclusion-be/src/core/entities"
	"github.com/educabot/alizia-inclusion-be/src/core/providers"
	"github.com/educabot/alizia-inclusion-be/src/observability"
)

// Cadencia del turno: NO es una personalidad distinta, solo cuánta urgencia/brevedad
// hay. Una sola Alizia (misma identidad y tono); lo único que cambia es el ritmo.
// Ortogonal a la Dimension (que dice QUÉ contexto traer, no CÓMO de breve responder).
// Ver alizia-comportamiento-flujo-v1.md §1.
const (
	CadenceInClass  = "assist" // en plena clase: breve, 1-3 acciones, al grano
	CadencePlanning = "guided" // planificando: puede tomarse un turno más para recopilar
)

type AssistClassroomRequest struct {
	OrgID          uuid.UUID
	UserID         int64
	ConversationID int64
	ClassroomID    int64
	StudentID      *int64
	Message        string
	// Mode es la cadencia (CadenceInClass / CadencePlanning), no una identidad.
	Mode string
	// Dimension (alumno / valija / tema) es la de la sesión abierta en /inclusion/open.
	// Opcional: si viene, dirige qué contexto trae el assembler; si no, se infiere del
	// alumno foco. Mismo vocabulario que OpenSession: un solo concepto de dimensión.
	Dimension string
	History   []providers.ChatMessage
}

func (r AssistClassroomRequest) Validate() error {
	if r.OrgID == uuid.Nil {
		return errOrgIDRequired
	}
	if r.Message == "" {
		return errMessageRequired
	}
	return nil
}

// ContentRef es una referencia liviana a un contenido pedagógico citado por el
// asistente (marker [CONTENT_ID:X] en el texto). El FE la usa para resolver el
// título del chip y deep-linkear.
type ContentRef struct {
	ID    int64  `json:"id"`
	Title string `json:"title"`
}

type AssistClassroomResponse struct {
	Response          string               `json:"response"`
	ConversationID    int64                `json:"conversation_id"`
	IdentifiedStudent *int64               `json:"identified_student,omitempty"`
	RecommendedDevice *int64               `json:"recommended_device,omitempty"`
	Adaptation        *GeneratedAdaptation `json:"adaptation,omitempty"`
	ReferencedContent []ContentRef         `json:"referenced_content,omitempty"`
	// Questions son las preguntas estructuradas (cajitas) que Alizia le hace al docente
	// este turno (marker [QUESTIONS_JSON]). El FE las renderiza como un Sheet con stepper.
	Questions []Question `json:"questions,omitempty"`
	// SourcesUsed expone, para este turno, qué tools de búsqueda ejecutó el modelo
	// (valija / alumno / RAG). Es la misma traza que el log chat.sources_used, pero
	// programática: el FE puede mostrar de dónde salió la respuesta y, sobre todo,
	// los evals pueden assertear que el modelo SÍ llamó las tools que correspondían.
	SourcesUsed SourcesUsed `json:"sources_used"`
}

// SourcesUsed es la vista pública (JSON) de sourcesSummary para el response del chat.
type SourcesUsed struct {
	Tools       []string `json:"tools,omitempty"`
	UsedValija  bool     `json:"used_valija"`
	UsedStudent bool     `json:"used_student"`
	UsedRAG     bool     `json:"used_rag"`
	StudentIDs  []int64  `json:"student_ids,omitempty"`
	RAGQueries  []string `json:"rag_queries,omitempty"`
	RAGHits     int      `json:"rag_hits"`
}

// newSourcesUsed proyecta el resumen interno al tipo público del response.
func newSourcesUsed(s sourcesSummary) SourcesUsed {
	return SourcesUsed{
		Tools:       s.Tools,
		UsedValija:  s.UsedValija,
		UsedStudent: s.UsedStudent,
		UsedRAG:     s.UsedRAG,
		StudentIDs:  s.StudentIDs,
		RAGQueries:  s.RAGQueries,
		RAGHits:     s.RAGHits,
	}
}

type AssistClassroom interface {
	Execute(ctx context.Context, req AssistClassroomRequest) (*AssistClassroomResponse, error)
}

// AssistClassroomDeps agrupa las dependencias del asistente de aula. Usar un struct
// (en vez de ~12 parámetros posicionales) mantiene legible el wiring y el constructor.
type AssistClassroomDeps struct {
	AI            providers.AIClient
	Students      providers.StudentProvider
	Profiles      providers.StudentProfileProvider
	Classrooms    providers.ClassroomProvider
	Devices       providers.DeviceProvider
	Conversations providers.ConversationProvider
	Summaries     providers.ConversationSummaryProvider
	Adaptations   providers.AdaptationProvider
	Content       providers.PedagogicalContentProvider
	Embedder      providers.Embedder
	RAG           providers.RAGSearchProvider
	Usage         providers.AIUsageProvider
	// CreateAdaptation / UpdateAdaptation persisten automáticamente el recurso que el
	// modelo genera en el chat (bloque ADAPTATION_JSON). Opcionales y nil-safe: si no
	// están inyectados, el chat funciona igual sin auto-guardar.
	CreateAdaptation CreateAdaptation
	UpdateAdaptation UpdateAdaptation
	// PromptCtx arma el contexto tipado del alumno/aula (Context Assembler, HU-2).
	// Opcional: si es nil o falla, el chat degrada al contexto base.
	PromptCtx BuildPromptContext
	Agentic   bool
}

type assistClassroomImpl struct {
	deps AssistClassroomDeps
}

func NewAssistClassroom(deps AssistClassroomDeps) AssistClassroom {
	return &assistClassroomImpl{deps: deps}
}

// resolveChatDimension elige la dimensión de contexto del turno: si el FE la manda
// explícita (la dimensión de la sesión abierta en /inclusion/open) la usa; si no, la
// infiere del alumno foco. Mode (cadencia) y dimensión quedan ortogonales: uno dice
// cuán breve responder, la otra qué contexto traer. Ver flujo §1.
func resolveChatDimension(dimension string, studentID *int64) string {
	switch d := strings.ToLower(strings.TrimSpace(dimension)); d {
	case DimensionStudent, DimensionToolkit, DimensionTopic:
		return d
	}
	if studentID != nil && *studentID > 0 {
		return DimensionStudent
	}
	return ""
}

// studentsDigest arma un resumen legible de los alumnos del aula para el log verbose.
func studentsDigest(students []entities.Student) string {
	parts := make([]string, len(students))
	for i := range students {
		s := &students[i]
		diff := ""
		if s.Profile != nil && len(s.Profile.Difficulties) > 0 {
			diff = " (" + strings.Join(s.Profile.Difficulties, ", ") + ")"
		}
		parts[i] = fmt.Sprintf("[%d]%s%s", s.ID, s.Name, diff)
	}
	return strings.Join(parts, "; ")
}

func (uc *assistClassroomImpl) Execute(ctx context.Context, req AssistClassroomRequest) (*AssistClassroomResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	// Correlación: todos los logs de este turno llevan org/user (y el request_id del ctx).
	ctx = observability.WithOrg(ctx, req.OrgID)
	ctx = observability.WithUser(ctx, req.UserID)

	// Contexto tipado del Context Assembler (HU-2): docente, alumno foco, diagnósticos,
	// PPI, adaptaciones previas y resúmenes. Opcional y nil-safe: si no está inyectado
	// o falla, el chat degrada al contexto base (devices + alumnos del aula).
	var pc *PromptContext
	if uc.deps.PromptCtx != nil && req.UserID > 0 {
		built, ctxErr := uc.deps.PromptCtx.Execute(ctx, BuildContextRequest{
			OrgID:     req.OrgID,
			UserID:    req.UserID,
			Dimension: resolveChatDimension(req.Dimension, req.StudentID),
			StudentID: req.StudentID,
		})
		if ctxErr != nil {
			slog.WarnContext(ctx, "chat.context_build_failed", "error", ctxErr)
		} else {
			pc = built
		}
	}

	var devices []entities.Device
	if pc != nil {
		devices = pc.DeviceCatalog
	} else {
		loaded, err := uc.deps.Devices.ListDevices(ctx, req.OrgID, nil)
		if err != nil {
			return nil, err
		}
		devices = loaded
	}

	// Todos los alumnos que el docente conoce en la org (no solo los del aula del turno):
	// así Alizia puede reconocer a un alumno aunque esté en otra aula.
	allStudents, _ := uc.deps.Students.List(ctx, req.OrgID)

	slog.InfoContext(ctx, "chat.context_loaded",
		"mode", req.Mode,
		"classroom_id", req.ClassroomID,
		"students_count", len(allStudents),
		"devices_count", len(devices),
		"context_assembled", pc != nil,
		observability.Text("students", studentsDigest(allStudents)),
	)

	// Recursos ya guardados en esta conversación: el modelo los ve en el prompt para
	// decidir si afina uno existente (devuelve su id → update) o crea uno nuevo. Best-effort.
	convResources := uc.loadConversationResources(ctx, req.OrgID, req.ConversationID)

	var systemPrompt string
	if req.Mode == CadencePlanning {
		systemPrompt = buildGuidedAssistPrompt(devices, allStudents, pc, convResources, uc.deps.Agentic)
	} else {
		systemPrompt = buildAssistSystemPrompt(devices, allStudents, pc, convResources, uc.deps.Agentic)
	}

	messages := make([]providers.ChatMessage, 0, len(req.History)+2)
	messages = append(messages, providers.ChatMessage{Role: "system", Content: systemPrompt})
	messages = append(messages, req.History...)
	messages = append(messages, providers.ChatMessage{Role: "user", Content: req.Message})
	messages = capMessages(messages, defaultMaxHistoryTokens)

	slog.InfoContext(ctx, "chat.prompt_built",
		"mode", req.Mode,
		"agentic", uc.deps.Agentic,
		"history_len", len(req.History),
		observability.Text("system_prompt", systemPrompt),
		observability.Text("user_message", req.Message),
	)

	var tools []providers.ToolDefinition
	if uc.deps.Agentic {
		tools = inclusionTools()
	}
	dispatcher := inclusionDispatcher{
		students:    uc.deps.Students,
		profiles:    uc.deps.Profiles,
		classrooms:  uc.deps.Classrooms,
		devices:     uc.deps.Devices,
		summaries:   uc.deps.Summaries,
		adaptations: uc.deps.Adaptations,
		content:     uc.deps.Content,
		embedder:    uc.deps.Embedder,
		rag:         uc.deps.RAG,
		userID:      req.UserID,
	}

	resp, trace, err := runAgenticChat(ctx, uc.deps.AI, messages, tools, dispatcher, req.OrgID, maxAgenticIterations, true)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", providers.ErrServiceUnavailable, err)
	}

	recordAIUsage(ctx, uc.deps.Usage, req.OrgID, req.UserID, "assist", resp.Usage)

	// Trazabilidad: de qué fuentes (valija / alumno / RAG) se sacó info este turno.
	sources := summarizeSources(trace)
	slog.InfoContext(ctx, "chat.sources_used",
		"used_valija", sources.UsedValija,
		"used_student", sources.UsedStudent,
		"student_ids", sources.StudentIDs,
		"used_rag", sources.UsedRAG,
		"rag_hits", sources.RAGHits,
		"tools", sources.Tools,
		observability.Text("rag_queries", strings.Join(sources.RAGQueries, " | ")),
	)

	slog.InfoContext(ctx, "chat.raw_response", "content", resp.Content)
	studentID := extractStudentID(resp.Content)
	deviceID := extractDeviceID(resp.Content)
	adaptation := extractAdaptationJSON(resp.Content)
	// Preguntas estructuradas (cajitas) que Alizia emite este turno. El bloque se quita
	// del texto visible en stripAdaptationBlock; el FE las renderiza desde este campo.
	questions := extractQuestions(resp.Content)
	if adaptation != nil && studentID != nil {
		adaptation.StudentID = studentID
	}
	// Contenido pedagógico citado este turno: el FE lo usa para resolver los chips
	// [CONTENT_ID:X]. Sale del back (lo que trajo el RAG), no de ids del modelo.
	referenced := contentRefsFromTrace(trace)
	// Quitamos solo el bloque ADAPTATION_JSON (ya extraído). Los markers
	// [STUDENT_ID:X]/[DEVICE_ID:X]/[CONTENT_ID:X] SÍ pasan: el FE los renderiza como
	// chips (nombre/título), nunca como id crudo.
	cleaned := sanitizeVisibleText(stripAdaptationBlock(resp.Content))

	// Guardrail duro de off-ramp: si el modelo cruzó el límite clínico (afirmar un
	// diagnóstico o dar una indicación clínica), reemplazamos la respuesta por una
	// derivación segura y descartamos la adaptación y las fuentes citadas (que ya no
	// aparecen en el texto). El prompt es la primera línea; esto es la red. Ver §5.
	if tripped, reason := crossedClinicalLine(cleaned); tripped {
		slog.WarnContext(ctx, "chat.guardrail_tripped",
			"reason", reason,
			observability.Text("original_response", cleaned),
		)
		cleaned = offRampMessage
		adaptation = nil
		referenced = nil
		questions = nil
	}

	convID, persistErr := uc.persistTurn(ctx, req, cleaned, studentID, deviceID, adaptation)
	if persistErr != nil {
		slog.WarnContext(ctx, "assist_classroom: persist turn failed", "error", persistErr, "user_id", req.UserID, "mode", req.Mode)
		convID = req.ConversationID
	}

	// Auto-persistencia del recurso: si el modelo emitió un ADAPTATION_JSON, lo creamos
	// (o actualizamos, si trae un id de esta conversación) en la tabla de recursos. Muta
	// adaptation.ID con el id persistido para devolvérselo al FE. Best-effort: no rompe el turno.
	if adaptation != nil {
		uc.persistAdaptation(ctx, req, convID, studentID, convResources, adaptation)
	}

	var idStudent, idDevice, totalTokens int64
	if studentID != nil {
		idStudent = *studentID
	}
	if deviceID != nil {
		idDevice = *deviceID
	}
	if resp.Usage != nil {
		totalTokens = int64(resp.Usage.TotalTokens)
	}
	slog.InfoContext(ctx, "chat.turn_done",
		"conversation_id", convID,
		"identified_student", idStudent,
		"recommended_device", idDevice,
		"has_adaptation", adaptation != nil,
		"questions_count", len(questions),
		"referenced_count", len(referenced),
		"total_tokens", totalTokens,
		observability.Text("response", cleaned),
	)

	return &AssistClassroomResponse{
		Response:          cleaned,
		ConversationID:    convID,
		IdentifiedStudent: studentID,
		RecommendedDevice: deviceID,
		Adaptation:        adaptation,
		ReferencedContent: referenced,
		Questions:         questions,
		SourcesUsed:       newSourcesUsed(sources),
	}, nil
}

func (uc *assistClassroomImpl) persistTurn(ctx context.Context, req AssistClassroomRequest, assistantContent string, studentID, deviceID *int64, adaptation *GeneratedAdaptation) (int64, error) {
	if uc.deps.Conversations == nil || req.UserID == 0 {
		return req.ConversationID, nil
	}
	mode := req.Mode
	if mode == "" {
		mode = CadenceInClass
	}
	metadata := map[string]any{}
	if studentID != nil {
		metadata["identified_student"] = *studentID
	}
	if deviceID != nil {
		metadata["recommended_device"] = *deviceID
	}
	if adaptation != nil {
		metadata["adaptation"] = adaptation
	}
	return uc.deps.Conversations.AppendTurn(ctx, providers.AppendTurnParams{
		ConversationID:   req.ConversationID,
		OrgID:            req.OrgID,
		UserID:           req.UserID,
		Mode:             mode,
		StudentID:        req.StudentID,
		UserContent:      req.Message,
		AssistantContent: assistantContent,
		Metadata:         metadata,
	})
}

// loadConversationResources lista los recursos ya guardados en la conversación, para
// inyectarlos en el prompt (el modelo decide afinar uno o crear otro). Best-effort:
// sin provider, sin conversación o ante error, devuelve nil y el chat sigue igual.
func (uc *assistClassroomImpl) loadConversationResources(ctx context.Context, orgID uuid.UUID, conversationID int64) []entities.Adaptation {
	if uc.deps.Adaptations == nil || conversationID == 0 {
		return nil
	}
	resources, err := uc.deps.Adaptations.List(ctx, providers.AdaptationFilter{
		OrgID:                orgID,
		SourceConversationID: &conversationID,
	})
	if err != nil {
		slog.WarnContext(ctx, "chat.conversation_resources_failed", "error", err, "conversation_id", conversationID)
		return nil
	}
	return resources
}

// persistAdaptation crea o actualiza el recurso pedagógico que el modelo generó este
// turno (bloque ADAPTATION_JSON). Si gen.ID viene y pertenece a un recurso de ESTA
// conversación, ACTUALIZA ese recurso; si no, CREA uno nuevo. Muta gen.ID con el id
// persistido para que el FE linkee al recurso real. Side-effect best-effort: nunca
// propaga error fatal (el docente ya tiene su respuesta).
func (uc *assistClassroomImpl) persistAdaptation(ctx context.Context, req AssistClassroomRequest, convID int64, studentID *int64, convResources []entities.Adaptation, gen *GeneratedAdaptation) {
	if gen == nil || req.UserID == 0 {
		return
	}

	// StudentID: prioridad al student_id del propio recurso (lo emite el modelo en el
	// ADAPTATION_JSON o lo copió el marcador [STUDENT_ID:X]); luego el marcador del texto;
	// y por último el alumno foco del request. Permite rellenar el alumno aunque el primer
	// guardado haya sido sin él (se crea un turno después).
	sid := gen.StudentID
	if sid == nil {
		sid = studentID
	}
	if sid == nil {
		sid = req.StudentID
	}

	// Update solo si el modelo devolvió un id que pertenece a esta conversación: así un id
	// inventado o de otra conversación no pisa un recurso ajeno (cae a create).
	if uc.deps.UpdateAdaptation != nil && gen.ID != nil && containsAdaptationID(convResources, *gen.ID) {
		strategy := gen.Strategy
		steps := entities.AdaptationSteps(gen.Steps)
		deviceIDs := gen.DeviceIDs
		upd := UpdateAdaptationRequest{
			OrgID:              req.OrgID,
			AdaptationID:       *gen.ID,
			Title:              &gen.Title,
			AdaptationStrategy: &strategy,
			Steps:              &steps,
			DeviceIDs:          &deviceIDs,
			RampID:             gen.RampID,
			// Backfill del alumno: si ya lo conocemos (sid != nil) lo rellenamos aunque el
			// recurso se haya creado sin alumno. Si sid == nil, Update no lo toca (no pisa).
			StudentID: sid,
			// Status NO se toca: si el docente ya lo marcó (probado/funcionó), se conserva.
		}
		if gen.Type != "" {
			atype := gen.Type
			upd.AdaptationType = &atype
		}
		updated, err := uc.deps.UpdateAdaptation.Execute(ctx, upd)
		if err != nil {
			slog.WarnContext(ctx, "chat.resource_save_failed", "error", err, "conversation_id", convID, "adaptation_id", *gen.ID)
			return
		}
		slog.InfoContext(ctx, "chat.resource_updated", "adaptation_id", updated.ID, "conversation_id", convID)
		gen.ID = &updated.ID
		return
	}

	if uc.deps.CreateAdaptation == nil {
		return
	}
	strategy := gen.Strategy
	created, err := uc.deps.CreateAdaptation.Execute(ctx, CreateAdaptationRequest{
		OrgID:                req.OrgID,
		StudentID:            sid,
		TeacherID:            req.UserID,
		Title:                gen.Title,
		AdaptationType:       gen.Type,
		AdaptationStrategy:   &strategy,
		Steps:                entities.AdaptationSteps(gen.Steps),
		RampID:               gen.RampID,
		DeviceIDs:            gen.DeviceIDs,
		SourceConversationID: &convID,
	})
	if err != nil {
		slog.WarnContext(ctx, "chat.resource_save_failed", "error", err, "conversation_id", convID)
		return
	}
	var sidLog int64
	if sid != nil {
		sidLog = *sid
	}
	slog.InfoContext(ctx, "chat.resource_saved",
		"adaptation_id", created.ID,
		"conversation_id", convID,
		"student_id", sidLog,
		observability.Text("title", created.Title),
	)
	gen.ID = &created.ID
}

// containsAdaptationID indica si id está en la lista de recursos de la conversación.
func containsAdaptationID(resources []entities.Adaptation, id int64) bool {
	for i := range resources {
		if resources[i].ID == id {
			return true
		}
	}
	return false
}
