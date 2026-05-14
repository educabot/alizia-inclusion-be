package inclusion

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/educabot/alizia-inclusion-be/src/core/providers"
)

type RecommendDeviceRequest struct {
	OrgID     uuid.UUID
	StudentID int64
	Subject   string
	Objective string
	Duration  string
	Dynamic   string
	Materials string
	History   []providers.ChatMessage
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
	if r.Objective == "" {
		return errObjectiveRequired
	}
	return nil
}

type RecommendDeviceResponse struct {
	Response string `json:"response"`
	DeviceID *int64 `json:"device_id,omitempty"`
}

type RecommendDevice interface {
	Execute(ctx context.Context, req RecommendDeviceRequest) (*RecommendDeviceResponse, error)
}

type recommendDeviceImpl struct {
	ai       providers.AIClient
	students providers.StudentProvider
	devices  providers.DeviceProvider
	ramps    providers.RampProvider
}

func NewRecommendDevice(ai providers.AIClient, students providers.StudentProvider, devices providers.DeviceProvider, ramps providers.RampProvider) RecommendDevice {
	return &recommendDeviceImpl{ai: ai, students: students, devices: devices, ramps: ramps}
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

	resp, err := uc.ai.Chat(ctx, messages)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", providers.ErrServiceUnavailable, err)
	}

	deviceID := extractDeviceID(resp.Content)

	return &RecommendDeviceResponse{
		Response: resp.Content,
		DeviceID: deviceID,
	}, nil
}
