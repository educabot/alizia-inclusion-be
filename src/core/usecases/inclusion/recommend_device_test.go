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

func TestRecommendDevice(t *testing.T) {
	ctx := context.Background()

	baseRequest := inclusion.RecommendDeviceRequest{
		OrgID:     testutil.TestOrgID,
		StudentID: 1,
		Subject:   "Matematicas",
		Objective: "Sumar fracciones",
	}

	setupMocks := func(aiResponse string, aiErr error) (*mocks.MockAIClient, *mocks.MockStudentProvider, *mocks.MockDeviceProvider, *mocks.MockRampProvider) {
		student := testutil.NewStudentWithProfile(1, 1, "Lucas", []string{"distraccion"})
		device := testutil.NewDevice(1, 1, "Timer Visual")
		return &mocks.MockAIClient{
				ChatFn: func(_ context.Context, _ []providers.ChatMessage) (*providers.ChatResponse, error) {
					if aiErr != nil {
						return nil, aiErr
					}
					return &providers.ChatResponse{Content: aiResponse}, nil
				},
			},
			&mocks.MockStudentProvider{
				GetStudentFn: func(_ context.Context, _ uuid.UUID, _ int64) (*entities.Student, error) {
					s := student
					return &s, nil
				},
			},
			&mocks.MockDeviceProvider{
				ListDevicesFn: func(_ context.Context, _ uuid.UUID, _ *int64) ([]entities.Device, error) {
					return []entities.Device{device}, nil
				},
			},
			&mocks.MockRampProvider{}
	}

	t.Run("returns recommendation with device", func(t *testing.T) {
		aiResp := `Recomiendo el Timer Visual [DEVICE_ID:1] para ayudar con la distraccion.
[ADAPTATION_JSON:{"title":"Timer para fracciones","type":"actividad_adaptada","strategy":"Usar timer","device_ids":[1],"device_names":["Timer Visual"]}]`
		ai, students, devices, ramps := setupMocks(aiResp, nil)

		got, err := inclusion.NewRecommendDevice(ai, students, devices, ramps).Execute(ctx, baseRequest)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if got.Response == "" {
			t.Error("expected non-empty response")
		}
		if got.DeviceID == nil || *got.DeviceID != 1 {
			t.Errorf("expected DeviceID 1, got %v", got.DeviceID)
		}
		if got.Adaptation == nil {
			t.Fatal("expected adaptation, got nil")
		}
		if got.Adaptation.Title != "Timer para fracciones" {
			t.Errorf("expected title %q, got %q", "Timer para fracciones", got.Adaptation.Title)
		}
	})

	t.Run("handles response without markers", func(t *testing.T) {
		ai, students, devices, ramps := setupMocks("Respuesta sin marcadores", nil)

		got, err := inclusion.NewRecommendDevice(ai, students, devices, ramps).Execute(ctx, baseRequest)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if got.DeviceID != nil {
			t.Errorf("expected nil DeviceID, got %d", *got.DeviceID)
		}
		if got.Adaptation != nil {
			t.Error("expected nil Adaptation")
		}
	})

	t.Run("wraps AI error as service unavailable", func(t *testing.T) {
		ai, students, devices, ramps := setupMocks("", errDB)

		_, err := inclusion.NewRecommendDevice(ai, students, devices, ramps).Execute(ctx, baseRequest)
		if err == nil {
			t.Fatal("expected error")
		}
		if !errors.Is(err, providers.ErrServiceUnavailable) {
			t.Errorf("expected ErrServiceUnavailable, got: %v", err)
		}
	})

	t.Run("rejects nil org_id", func(t *testing.T) {
		ai, students, devices, ramps := setupMocks("", nil)
		req := baseRequest
		req.OrgID = uuid.Nil
		_, err := inclusion.NewRecommendDevice(ai, students, devices, ramps).Execute(ctx, req)
		if !errors.Is(err, providers.ErrValidation) {
			t.Errorf("expected ErrValidation, got: %v", err)
		}
	})

	t.Run("rejects zero student_id", func(t *testing.T) {
		ai, students, devices, ramps := setupMocks("", nil)
		req := baseRequest
		req.StudentID = 0
		_, err := inclusion.NewRecommendDevice(ai, students, devices, ramps).Execute(ctx, req)
		if !errors.Is(err, providers.ErrValidation) {
			t.Errorf("expected ErrValidation, got: %v", err)
		}
	})

	t.Run("rejects empty subject", func(t *testing.T) {
		ai, students, devices, ramps := setupMocks("", nil)
		req := baseRequest
		req.Subject = ""
		_, err := inclusion.NewRecommendDevice(ai, students, devices, ramps).Execute(ctx, req)
		if !errors.Is(err, providers.ErrValidation) {
			t.Errorf("expected ErrValidation, got: %v", err)
		}
	})

}
