package inclusion_test

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"

	"github.com/educabot/alizia-inclusion-be/src/core/entities"
	"github.com/educabot/alizia-inclusion-be/src/core/providers"
	"github.com/educabot/alizia-inclusion-be/src/core/providers/mocks"
	"github.com/educabot/alizia-inclusion-be/src/core/usecases/inclusion"
	"github.com/educabot/alizia-inclusion-be/src/testutil"
)

func assistClassroomSetupMocks(aiResponse string, aiErr error) (*mocks.MockAIClient, *mocks.MockStudentProvider, *mocks.MockDeviceProvider, *mocks.MockConversationProvider, *mocks.MockAIUsageProvider) {
	return &mocks.MockAIClient{
			ChatFn: func(_ context.Context, _ []providers.ChatMessage) (*providers.ChatResponse, error) {
				if aiErr != nil {
					return nil, aiErr
				}
				return &providers.ChatResponse{Content: aiResponse}, nil
			},
		},
		&mocks.MockStudentProvider{
			ListByClassroomFn: func(_ context.Context, _ uuid.UUID, _ int64) ([]entities.Student, error) {
				s := testutil.NewStudent(1, 1, "Lucas")
				return []entities.Student{s}, nil
			},
		},
		&mocks.MockDeviceProvider{
			ListDevicesFn: func(_ context.Context, _ uuid.UUID, _ *int64) ([]entities.Device, error) {
				d := testutil.NewDevice(1, 1, "Pictogramas")
				return []entities.Device{d}, nil
			},
		},
		&mocks.MockConversationProvider{},
		&mocks.MockAIUsageProvider{}
}

var assistClassroomBaseRequest = inclusion.AssistClassroomRequest{
	OrgID:       testutil.TestOrgID,
	ClassroomID: 1,
	Message:     "Lucas no puede concentrarse",
}

func TestAssistClassroom_ReturnsAssistResponse(t *testing.T) {
	ctx := context.Background()

	ai, students, devices, conversations, usage := assistClassroomSetupMocks("Podrias usar pictogramas [DEVICE_ID:1]", nil)

	got, err := inclusion.NewAssistClassroom(ai, students, devices, conversations, usage, false).Execute(ctx, assistClassroomBaseRequest)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.Response == "" {
		t.Error("expected non-empty response")
	}
}

func TestAssistClassroom_WorksInGuidedMode(t *testing.T) {
	ctx := context.Background()

	ai, students, devices, conversations, usage := assistClassroomSetupMocks("Para que alumno necesitas la adaptacion?", nil)
	req := assistClassroomBaseRequest
	req.Mode = "guided"

	got, err := inclusion.NewAssistClassroom(ai, students, devices, conversations, usage, false).Execute(ctx, req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.Response == "" {
		t.Error("expected non-empty response")
	}
}

func TestAssistClassroom_WrapsAIError(t *testing.T) {
	ctx := context.Background()

	ai, students, devices, conversations, usage := assistClassroomSetupMocks("", errDB)

	_, err := inclusion.NewAssistClassroom(ai, students, devices, conversations, usage, false).Execute(ctx, assistClassroomBaseRequest)
	if err == nil {
		t.Fatal("expected error")
	}
	if !errors.Is(err, providers.ErrServiceUnavailable) {
		t.Errorf("expected ErrServiceUnavailable, got: %v", err)
	}
}

func TestAssistClassroom_RejectsNilOrgID(t *testing.T) {
	ctx := context.Background()

	ai, students, devices, conversations, usage := assistClassroomSetupMocks("", nil)
	req := assistClassroomBaseRequest
	req.OrgID = uuid.Nil

	_, err := inclusion.NewAssistClassroom(ai, students, devices, conversations, usage, false).Execute(ctx, req)
	if !errors.Is(err, providers.ErrValidation) {
		t.Errorf("expected ErrValidation, got: %v", err)
	}
}

func TestAssistClassroom_RejectsEmptyMessage(t *testing.T) {
	ctx := context.Background()

	ai, students, devices, conversations, usage := assistClassroomSetupMocks("", nil)
	req := assistClassroomBaseRequest
	req.Message = ""

	_, err := inclusion.NewAssistClassroom(ai, students, devices, conversations, usage, false).Execute(ctx, req)
	if !errors.Is(err, providers.ErrValidation) {
		t.Errorf("expected ErrValidation, got: %v", err)
	}
}

func TestAssistClassroom_PersistsTurnWhenUserIDPresent(t *testing.T) {
	ctx := context.Background()

	ai, students, devices, conversations, usage := assistClassroomSetupMocks("Podrias usar pictogramas [DEVICE_ID:1]", nil)
	var captured providers.AppendTurnParams
	conversations.AppendTurnFn = func(_ context.Context, p providers.AppendTurnParams) (int64, error) {
		captured = p
		return 42, nil
	}
	req := assistClassroomBaseRequest
	req.UserID = 7

	got, err := inclusion.NewAssistClassroom(ai, students, devices, conversations, usage, false).Execute(ctx, req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.ConversationID != 42 {
		t.Errorf("expected conversation_id 42, got %d", got.ConversationID)
	}
	if captured.UserID != 7 {
		t.Errorf("expected captured UserID 7, got %d", captured.UserID)
	}
	if captured.UserContent != assistClassroomBaseRequest.Message {
		t.Errorf("expected UserContent %q, got %q", assistClassroomBaseRequest.Message, captured.UserContent)
	}
}

func TestAssistClassroom_SkipsPersistenceWhenUserIDMissing(t *testing.T) {
	ctx := context.Background()

	ai, students, devices, conversations, usage := assistClassroomSetupMocks("ok", nil)
	called := false
	conversations.AppendTurnFn = func(_ context.Context, _ providers.AppendTurnParams) (int64, error) {
		called = true
		return 99, nil
	}

	_, err := inclusion.NewAssistClassroom(ai, students, devices, conversations, usage, false).Execute(ctx, assistClassroomBaseRequest)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if called {
		t.Error("AppendTurn should not be invoked without UserID")
	}
}

func TestAssistClassroom_RecordsTokenUsageWhenPresent(t *testing.T) {
	ctx := context.Background()

	ai, students, devices, conversations, usage := assistClassroomSetupMocks("ok", nil)
	ai.ChatFn = func(_ context.Context, _ []providers.ChatMessage) (*providers.ChatResponse, error) {
		return &providers.ChatResponse{
			Content: "ok",
			Usage:   &providers.TokenUsage{PromptTokens: 10, CompletionTokens: 5, TotalTokens: 15},
		}, nil
	}
	var captured providers.AIUsageRecord
	usage.RecordFn = func(_ context.Context, r providers.AIUsageRecord) error {
		captured = r
		return nil
	}
	req := assistClassroomBaseRequest
	req.UserID = 7

	_, err := inclusion.NewAssistClassroom(ai, students, devices, conversations, usage, false).Execute(ctx, req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if captured.TotalTokens != 15 {
		t.Errorf("expected total_tokens 15, got %d", captured.TotalTokens)
	}
	if captured.Mode != "assist" {
		t.Errorf("expected mode 'assist', got %q", captured.Mode)
	}
}
