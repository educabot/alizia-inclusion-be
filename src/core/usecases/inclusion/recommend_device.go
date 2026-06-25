package inclusion

import (
	"context"
	"log/slog"
	"time"

	"github.com/google/uuid"

	"github.com/educabot/alizia-inclusion-be/src/core/providers"
	"github.com/educabot/alizia-inclusion-be/src/core/usecases/inclusion/prompts"
)

type RecommendDeviceRequest struct {
	OrgID          uuid.UUID
	UserID         int64
	ConversationID int64
	StudentID      int64
	Subject        string
	Objective      string
	Duration       string
	Dynamic        string
	Materials      string
	History        []providers.ChatMessage
}

func (r RecommendDeviceRequest) Validate() error {
	if r.OrgID == uuid.Nil {
		return errOrgIDRequired
	}
	if r.StudentID <= 0 {
		return errStudentIDRequired
	}
	if r.Subject == "" {
		return errSubjectRequired
	}
	return nil
}

type RecommendDeviceResponse struct {
	Response       string               `json:"response"`
	ConversationID int64                `json:"conversation_id"`
	DeviceID       *int64               `json:"device_id,omitempty"`
	Adaptation     *GeneratedAdaptation `json:"adaptation,omitempty"`
}

type RecommendDevice interface {
	Execute(ctx context.Context, req RecommendDeviceRequest) (*RecommendDeviceResponse, error)
}

type recommendDeviceImpl struct {
	ai            providers.AIClient
	students      providers.StudentProvider
	devices       providers.DeviceProvider
	ramps         providers.RampProvider
	conversations providers.ConversationProvider
	usage         providers.AIUsageProvider
}

func NewRecommendDevice(ai providers.AIClient, students providers.StudentProvider, devices providers.DeviceProvider, ramps providers.RampProvider, conversations providers.ConversationProvider, usage providers.AIUsageProvider) RecommendDevice {
	return &recommendDeviceImpl{ai: ai, students: students, devices: devices, ramps: ramps, conversations: conversations, usage: usage}
}

func (uc *recommendDeviceImpl) Execute(ctx context.Context, req RecommendDeviceRequest) (*RecommendDeviceResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	student, err := uc.students.GetStudent(ctx, req.OrgID, req.StudentID)
	if err != nil {
		return nil, err
	}

	devices, err := uc.devices.ListDevices(ctx, req.OrgID, nil)
	if err != nil {
		return nil, err
	}

	systemPrompt := prompts.RecommendSystem(devices)
	userPrompt := buildRecommendUserPrompt(student, req)
	messages := buildChatMessages(systemPrompt, req.History, userPrompt)

	start := time.Now()
	resp, err := uc.ai.Chat(ctx, messages)
	if err != nil {
		return nil, wrapServiceUnavailable(err)
	}
	latencyMs := int(time.Since(start).Milliseconds())

	// Guardrail: recommend is the heaviest path for DEVICE_ID/ADAPTATION_JSON;
	// a hallucinated device ID must never reach the teacher. Fall back to off-ramp, same as assist.
	guardAnswer(ctx, resp, devices, "usecase", "recommend_device", "user_id", req.UserID, "student_id", req.StudentID)

	deviceID := extractDeviceID(resp.Content)
	adaptation := extractAdaptationJSON(resp.Content)

	convID, persistErr := uc.persistTurn(ctx, req, userPrompt, resp.Content, deviceID, adaptation)
	if persistErr != nil {
		slog.WarnContext(ctx, "recommend_device: persist turn failed", "error", persistErr, "user_id", req.UserID, "student_id", req.StudentID)
		convID = req.ConversationID
	}

	// Per-turn trace: IDs only, no PII. Best-effort.
	snapshot := map[string]any{"student_id": req.StudentID}
	if deviceID != nil {
		snapshot["recommended_device_id"] = *deviceID
	}
	recordAIUsage(ctx, uc.usage, aiTrace{
		orgID: req.OrgID, userID: req.UserID, mode: modeRecommend,
		model: uc.ai.Model(), latencyMs: latencyMs,
		conversationID: convID, usage: resp.Usage, context: snapshot,
	})

	return &RecommendDeviceResponse{
		Response:       resp.Content,
		ConversationID: convID,
		DeviceID:       deviceID,
		Adaptation:     adaptation,
	}, nil
}

func (uc *recommendDeviceImpl) persistTurn(ctx context.Context, req RecommendDeviceRequest, userContent, assistantContent string, deviceID *int64, adaptation *GeneratedAdaptation) (int64, error) {
	if uc.conversations == nil || req.UserID == 0 {
		return req.ConversationID, nil
	}
	metadata := map[string]any{
		metaKeySubject: req.Subject,
	}
	if deviceID != nil {
		metadata[metaKeyRecommendedDevice] = *deviceID
	}
	if adaptation != nil {
		metadata[metaKeyAdaptation] = adaptation
	}
	studentIDCopy := req.StudentID
	return uc.conversations.AppendTurn(ctx, providers.AppendTurnParams{
		ConversationID:   req.ConversationID,
		OrgID:            req.OrgID,
		UserID:           req.UserID,
		Mode:             modeRecommend,
		StudentID:        &studentIDCopy,
		UserContent:      userContent,
		AssistantContent: assistantContent,
		Metadata:         metadata,
	})
}
