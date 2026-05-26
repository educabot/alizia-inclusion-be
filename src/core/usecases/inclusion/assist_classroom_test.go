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

func TestAssistClassroom(t *testing.T) {
	ctx := context.Background()

	baseRequest := inclusion.AssistClassroomRequest{
		OrgID:       testutil.TestOrgID,
		ClassroomID: 1,
		Message:     "Lucas no puede concentrarse",
	}

	setupMocks := func(aiResponse string, aiErr error) (*mocks.MockAIClient, *mocks.MockStudentProvider, *mocks.MockDeviceProvider, *mocks.MockConversationProvider, *mocks.MockAIUsageProvider) {
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

	t.Run("returns assist response", func(t *testing.T) {
		ai, students, devices, conversations, usage := setupMocks("Podrias usar pictogramas [DEVICE_ID:1]", nil)

		got, err := inclusion.NewAssistClassroom(ai, students, devices, conversations, usage).Execute(ctx, baseRequest)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if got.Response == "" {
			t.Error("expected non-empty response")
		}
	})

	t.Run("works in guided mode", func(t *testing.T) {
		ai, students, devices, conversations, usage := setupMocks("Para que alumno necesitas la adaptacion?", nil)
		req := baseRequest
		req.Mode = "guided"

		got, err := inclusion.NewAssistClassroom(ai, students, devices, conversations, usage).Execute(ctx, req)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if got.Response == "" {
			t.Error("expected non-empty response")
		}
	})

	t.Run("wraps AI error", func(t *testing.T) {
		ai, students, devices, conversations, usage := setupMocks("", errDB)

		_, err := inclusion.NewAssistClassroom(ai, students, devices, conversations, usage).Execute(ctx, baseRequest)
		if err == nil {
			t.Fatal("expected error")
		}
		if !errors.Is(err, providers.ErrServiceUnavailable) {
			t.Errorf("expected ErrServiceUnavailable, got: %v", err)
		}
	})

	t.Run("rejects nil org_id", func(t *testing.T) {
		ai, students, devices, conversations, usage := setupMocks("", nil)
		req := baseRequest
		req.OrgID = uuid.Nil
		_, err := inclusion.NewAssistClassroom(ai, students, devices, conversations, usage).Execute(ctx, req)
		if !errors.Is(err, providers.ErrValidation) {
			t.Errorf("expected ErrValidation, got: %v", err)
		}
	})

	t.Run("rejects empty message", func(t *testing.T) {
		ai, students, devices, conversations, usage := setupMocks("", nil)
		req := baseRequest
		req.Message = ""
		_, err := inclusion.NewAssistClassroom(ai, students, devices, conversations, usage).Execute(ctx, req)
		if !errors.Is(err, providers.ErrValidation) {
			t.Errorf("expected ErrValidation, got: %v", err)
		}
	})

	t.Run("persists turn when user_id present", func(t *testing.T) {
		ai, students, devices, conversations, usage := setupMocks("Podrias usar pictogramas [DEVICE_ID:1]", nil)
		var captured providers.AppendTurnParams
		conversations.AppendTurnFn = func(_ context.Context, p providers.AppendTurnParams) (int64, error) {
			captured = p
			return 42, nil
		}
		req := baseRequest
		req.UserID = 7

		got, err := inclusion.NewAssistClassroom(ai, students, devices, conversations, usage).Execute(ctx, req)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if got.ConversationID != 42 {
			t.Errorf("expected conversation_id 42, got %d", got.ConversationID)
		}
		if captured.UserID != 7 {
			t.Errorf("expected captured UserID 7, got %d", captured.UserID)
		}
		if captured.UserContent != baseRequest.Message {
			t.Errorf("expected UserContent %q, got %q", baseRequest.Message, captured.UserContent)
		}
	})

	t.Run("skips persistence when user_id missing", func(t *testing.T) {
		ai, students, devices, conversations, usage := setupMocks("ok", nil)
		called := false
		conversations.AppendTurnFn = func(_ context.Context, _ providers.AppendTurnParams) (int64, error) {
			called = true
			return 99, nil
		}

		_, err := inclusion.NewAssistClassroom(ai, students, devices, conversations, usage).Execute(ctx, baseRequest)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if called {
			t.Error("AppendTurn should not be invoked without UserID")
		}
	})

	t.Run("records token usage when present", func(t *testing.T) {
		ai, students, devices, conversations, usage := setupMocks("ok", nil)
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
		req := baseRequest
		req.UserID = 7

		_, err := inclusion.NewAssistClassroom(ai, students, devices, conversations, usage).Execute(ctx, req)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if captured.TotalTokens != 15 {
			t.Errorf("expected total_tokens 15, got %d", captured.TotalTokens)
		}
		if captured.Mode != "assist" {
			t.Errorf("expected mode 'assist', got %q", captured.Mode)
		}
	})
}
