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

	setupMocks := func(aiResponse string, aiErr error) (*mocks.MockAIClient, *mocks.MockStudentProvider, *mocks.MockDeviceProvider) {
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
			}
	}

	t.Run("returns assist response", func(t *testing.T) {
		ai, students, devices := setupMocks("Podrias usar pictogramas [DEVICE_ID:1]", nil)

		got, err := inclusion.NewAssistClassroom(ai, students, devices).Execute(ctx, baseRequest)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if got.Response == "" {
			t.Error("expected non-empty response")
		}
	})

	t.Run("works in guided mode", func(t *testing.T) {
		ai, students, devices := setupMocks("Para que alumno necesitas la adaptacion?", nil)
		req := baseRequest
		req.Mode = "guided"

		got, err := inclusion.NewAssistClassroom(ai, students, devices).Execute(ctx, req)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if got.Response == "" {
			t.Error("expected non-empty response")
		}
	})

	t.Run("wraps AI error", func(t *testing.T) {
		ai, students, devices := setupMocks("", errDB)

		_, err := inclusion.NewAssistClassroom(ai, students, devices).Execute(ctx, baseRequest)
		if err == nil {
			t.Fatal("expected error")
		}
		if !errors.Is(err, providers.ErrServiceUnavailable) {
			t.Errorf("expected ErrServiceUnavailable, got: %v", err)
		}
	})

	t.Run("rejects nil org_id", func(t *testing.T) {
		ai, students, devices := setupMocks("", nil)
		req := baseRequest
		req.OrgID = uuid.Nil
		_, err := inclusion.NewAssistClassroom(ai, students, devices).Execute(ctx, req)
		if !errors.Is(err, providers.ErrValidation) {
			t.Errorf("expected ErrValidation, got: %v", err)
		}
	})

	t.Run("rejects empty message", func(t *testing.T) {
		ai, students, devices := setupMocks("", nil)
		req := baseRequest
		req.Message = ""
		_, err := inclusion.NewAssistClassroom(ai, students, devices).Execute(ctx, req)
		if !errors.Is(err, providers.ErrValidation) {
			t.Errorf("expected ErrValidation, got: %v", err)
		}
	})
}
