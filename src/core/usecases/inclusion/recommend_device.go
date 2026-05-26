package inclusion

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/google/uuid"

	"github.com/educabot/alizia-inclusion-be/src/core/providers"
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

	systemPrompt := buildRecommendSystemPrompt(devices)
	userPrompt := buildRecommendUserPrompt(student, req)

	messages := make([]providers.ChatMessage, 0, len(req.History)+2)
	messages = append(messages, providers.ChatMessage{Role: "system", Content: systemPrompt})
	messages = append(messages, req.History...)
	messages = append(messages, providers.ChatMessage{Role: "user", Content: userPrompt})
	messages = capMessages(messages, defaultMaxHistoryTokens)

	resp, err := uc.ai.Chat(ctx, messages)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", providers.ErrServiceUnavailable, err)
	}

	recordAIUsage(ctx, uc.usage, req.OrgID, req.UserID, "recommend", resp.Usage)

	deviceID := extractDeviceID(resp.Content)
	adaptation := extractAdaptationJSON(resp.Content)

	convID, persistErr := uc.persistTurn(ctx, req, userPrompt, resp.Content, deviceID, adaptation)
	if persistErr != nil {
		slog.WarnContext(ctx, "recommend_device: persist turn failed", "error", persistErr, "user_id", req.UserID, "student_id", req.StudentID)
		convID = req.ConversationID
	}

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
		"subject": req.Subject,
	}
	if deviceID != nil {
		metadata["recommended_device"] = *deviceID
	}
	if adaptation != nil {
		metadata["adaptation"] = adaptation
	}
	studentIDCopy := req.StudentID
	return uc.conversations.AppendTurn(ctx, providers.AppendTurnParams{
		ConversationID:   req.ConversationID,
		OrgID:            req.OrgID,
		UserID:           req.UserID,
		Mode:             "recommend",
		StudentID:        &studentIDCopy,
		UserContent:      userContent,
		AssistantContent: assistantContent,
		Metadata:         metadata,
	})
}
