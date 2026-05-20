package inclusion

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/educabot/alizia-inclusion-be/src/core/providers"
)

type AssistClassroomRequest struct {
	OrgID       uuid.UUID
	ClassroomID int64
	StudentID   *int64
	Message     string
	Mode        string
	History     []providers.ChatMessage
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
	IdentifiedStudent *int64               `json:"identified_student,omitempty"`
	RecommendedDevice *int64               `json:"recommended_device,omitempty"`
	Adaptation        *GeneratedAdaptation `json:"adaptation,omitempty"`
}

type AssistClassroom interface {
	Execute(ctx context.Context, req AssistClassroomRequest) (*AssistClassroomResponse, error)
}

type assistClassroomImpl struct {
	ai       providers.AIClient
	students providers.StudentProvider
	devices  providers.DeviceProvider
}

func NewAssistClassroom(ai providers.AIClient, students providers.StudentProvider, devices providers.DeviceProvider) AssistClassroom {
	return &assistClassroomImpl{ai: ai, students: students, devices: devices}
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

	var systemPrompt string
	if req.Mode == "guided" {
		systemPrompt = buildGuidedAssistPrompt(devices, allStudents)
	} else {
		systemPrompt = buildAssistSystemPrompt(devices, allStudents)
	}

	messages := make([]providers.ChatMessage, 0, len(req.History)+2)
	messages = append(messages, providers.ChatMessage{Role: "system", Content: systemPrompt})
	messages = append(messages, req.History...)
	messages = append(messages, providers.ChatMessage{Role: "user", Content: req.Message})

	resp, err := uc.ai.Chat(ctx, messages)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", providers.ErrServiceUnavailable, err)
	}

	studentID := extractStudentID(resp.Content)
	deviceID := extractDeviceID(resp.Content)
	adaptation := extractAdaptationJSON(resp.Content)

	return &AssistClassroomResponse{
		Response:          resp.Content,
		IdentifiedStudent: studentID,
		RecommendedDevice: deviceID,
		Adaptation:        adaptation,
	}, nil
}
