package inclusion

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/google/uuid"

	"github.com/educabot/alizia-inclusion-be/src/core/providers"
	"github.com/educabot/alizia-inclusion-be/src/core/usecases/inclusion/prompts"
)

type AssistClassroomRequest struct {
	OrgID          uuid.UUID
	UserID         int64
	ConversationID int64
	ClassroomID    int64
	StudentID      *int64
	Message        string
	Mode           string
	History        []providers.ChatMessage
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

type AssistClassroomResponse struct {
	Response          string               `json:"response"`
	ConversationID    int64                `json:"conversation_id"`
	IdentifiedStudent *int64               `json:"identified_student,omitempty"`
	RecommendedDevice *int64               `json:"recommended_device,omitempty"`
	Adaptation        *GeneratedAdaptation `json:"adaptation,omitempty"`
}

type AssistClassroom interface {
	Execute(ctx context.Context, req AssistClassroomRequest) (*AssistClassroomResponse, error)
}

type assistClassroomImpl struct {
	ai            providers.AIClient
	students      providers.StudentProvider
	devices       providers.DeviceProvider
	conversations providers.ConversationProvider
	summaries     providers.ConversationSummaryProvider
	adaptations   providers.AdaptationProvider
	content       providers.PedagogicalContentProvider
	usage         providers.AIUsageProvider
	agentic       bool
}

func NewAssistClassroom(ai providers.AIClient, students providers.StudentProvider, devices providers.DeviceProvider, conversations providers.ConversationProvider, summaries providers.ConversationSummaryProvider, adaptations providers.AdaptationProvider, content providers.PedagogicalContentProvider, usage providers.AIUsageProvider, agentic bool) AssistClassroom {
	return &assistClassroomImpl{ai: ai, students: students, devices: devices, conversations: conversations, summaries: summaries, adaptations: adaptations, content: content, usage: usage, agentic: agentic}
}

func (uc *assistClassroomImpl) Execute(ctx context.Context, req AssistClassroomRequest) (*AssistClassroomResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	devices, err := uc.devices.ListDevices(ctx, req.OrgID, nil)
	if err != nil {
		return nil, err
	}

	allStudents, _ := uc.students.ListByClassroom(ctx, req.OrgID, req.ClassroomID)

	// Un solo modo (HU-6, T-6.1): el marco pedagógico vive en el paquete prompts.
	systemPrompt := prompts.AssistSystem(devices, allStudents)

	messages := make([]providers.ChatMessage, 0, len(req.History)+2)
	messages = append(messages, providers.ChatMessage{Role: "system", Content: systemPrompt})
	messages = append(messages, req.History...)
	messages = append(messages, providers.ChatMessage{Role: "user", Content: req.Message})
	messages = capMessages(messages, defaultMaxHistoryTokens)

	var tools []providers.ToolDefinition
	if uc.agentic {
		tools = inclusionTools()
	}
	dispatcher := inclusionDispatcher{students: uc.students, devices: uc.devices, summaries: uc.summaries, adaptations: uc.adaptations, content: uc.content, userID: req.UserID, conversationID: req.ConversationID}

	start := time.Now()
	resp, toolCalls, err := runAgenticChat(ctx, uc.ai, messages, tools, dispatcher, req.OrgID, maxAgenticIterations)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", providers.ErrServiceUnavailable, err)
	}
	latencyMs := int(time.Since(start).Milliseconds())

	// Guardrail por código (HU-6): si la respuesta cruza un límite duro (ej. cita
	// un DEVICE_ID inexistente) no se la mostramos al docente; caemos al off-ramp.
	if gr := validateAnswer(resp.Content, deviceCatalogSet(devices)); !gr.Valid {
		slog.WarnContext(ctx, "assist_classroom: guardrail rejected answer",
			"violations", gr.Violations, "user_id", req.UserID, "mode", req.Mode)
		resp.Content = prompts.OffRampInvalidOutput
	}

	studentID := extractStudentID(resp.Content)
	deviceID := extractDeviceID(resp.Content)
	adaptation := extractAdaptationJSON(resp.Content)

	convID, persistErr := uc.persistTurn(ctx, req, resp.Content, studentID, deviceID, adaptation)
	if persistErr != nil {
		slog.WarnContext(ctx, "assist_classroom: persist turn failed", "error", persistErr, "user_id", req.UserID, "mode", req.Mode)
		convID = req.ConversationID
	}

	// Traza por turno (HU-6, T-6.5): solo IDs, sin PII. Best-effort.
	recordAIUsage(ctx, uc.usage, aiTrace{
		orgID: req.OrgID, userID: req.UserID, mode: "assist",
		model: uc.ai.Model(), latencyMs: latencyMs, toolCalls: toolCalls,
		conversationID: convID, usage: resp.Usage,
		context: assistContextSnapshot(req, studentID, deviceID),
	})

	return &AssistClassroomResponse{
		Response:          resp.Content,
		ConversationID:    convID,
		IdentifiedStudent: studentID,
		RecommendedDevice: deviceID,
		Adaptation:        adaptation,
	}, nil
}

// assistContextSnapshot arma el snapshot de contexto del turno SOLO con IDs (sin
// PII): qué aula/alumno se mencionó y qué dispositivo se recomendó (HU-6, T-6.5).
func assistContextSnapshot(req AssistClassroomRequest, studentID, deviceID *int64) map[string]any {
	snap := map[string]any{}
	if req.ClassroomID > 0 {
		snap["classroom_id"] = req.ClassroomID
	}
	if req.StudentID != nil {
		snap["student_id"] = *req.StudentID
	}
	if studentID != nil {
		snap["identified_student_id"] = *studentID
	}
	if deviceID != nil {
		snap["recommended_device_id"] = *deviceID
	}
	return snap
}

func (uc *assistClassroomImpl) persistTurn(ctx context.Context, req AssistClassroomRequest, assistantContent string, studentID, deviceID *int64, adaptation *GeneratedAdaptation) (int64, error) {
	if uc.conversations == nil || req.UserID == 0 {
		return req.ConversationID, nil
	}
	mode := req.Mode
	if mode == "" {
		mode = "assist"
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
	return uc.conversations.AppendTurn(ctx, providers.AppendTurnParams{
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
