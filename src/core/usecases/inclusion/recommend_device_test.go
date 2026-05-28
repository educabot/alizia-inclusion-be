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

func newRecommendMocks(aiResponse string, aiErr error) (*mocks.MockAIClient, *mocks.MockStudentProvider, *mocks.MockDeviceProvider, *mocks.MockRampProvider, *mocks.MockConversationProvider, *mocks.MockAIUsageProvider) {
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
		&mocks.MockRampProvider{},
		&mocks.MockConversationProvider{},
		&mocks.MockAIUsageProvider{}
}

var baseRecommendRequest = inclusion.RecommendDeviceRequest{
	OrgID:     testutil.TestOrgID,
	StudentID: 1,
	Subject:   "Matematicas",
	Objective: "Sumar fracciones",
}

func TestRecommendDevice_ReturnsRecommendationWithDevice(t *testing.T) {
	ctx := context.Background()
	aiResp := `Recomiendo el Timer Visual [DEVICE_ID:1] para ayudar con la distraccion.
[ADAPTATION_JSON:{"title":"Timer para fracciones","type":"actividad_adaptada","strategy":"Usar timer","device_ids":[1],"device_names":["Timer Visual"]}]`
	ai, students, devices, ramps, conversations, usage := newRecommendMocks(aiResp, nil)

	got, err := inclusion.NewRecommendDevice(ai, students, devices, ramps, conversations, usage).Execute(ctx, baseRecommendRequest)

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
}

func TestRecommendDevice_HandlesResponseWithoutMarkers(t *testing.T) {
	ctx := context.Background()
	ai, students, devices, ramps, conversations, usage := newRecommendMocks("Respuesta sin marcadores", nil)

	got, err := inclusion.NewRecommendDevice(ai, students, devices, ramps, conversations, usage).Execute(ctx, baseRecommendRequest)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.DeviceID != nil {
		t.Errorf("expected nil DeviceID, got %d", *got.DeviceID)
	}
	if got.Adaptation != nil {
		t.Error("expected nil Adaptation")
	}
}

func TestRecommendDevice_WrapsAIErrorAsServiceUnavailable(t *testing.T) {
	ctx := context.Background()
	ai, students, devices, ramps, conversations, usage := newRecommendMocks("", errDB)

	_, err := inclusion.NewRecommendDevice(ai, students, devices, ramps, conversations, usage).Execute(ctx, baseRecommendRequest)

	if err == nil {
		t.Fatal("expected error")
	}
	if !errors.Is(err, providers.ErrServiceUnavailable) {
		t.Errorf("expected ErrServiceUnavailable, got: %v", err)
	}
}

func TestRecommendDevice_RejectsNilOrgID(t *testing.T) {
	ctx := context.Background()
	ai, students, devices, ramps, conversations, usage := newRecommendMocks("", nil)
	req := baseRecommendRequest
	req.OrgID = uuid.Nil

	_, err := inclusion.NewRecommendDevice(ai, students, devices, ramps, conversations, usage).Execute(ctx, req)

	if !errors.Is(err, providers.ErrValidation) {
		t.Errorf("expected ErrValidation, got: %v", err)
	}
}

func TestRecommendDevice_RejectsZeroStudentID(t *testing.T) {
	ctx := context.Background()
	ai, students, devices, ramps, conversations, usage := newRecommendMocks("", nil)
	req := baseRecommendRequest
	req.StudentID = 0

	_, err := inclusion.NewRecommendDevice(ai, students, devices, ramps, conversations, usage).Execute(ctx, req)

	if !errors.Is(err, providers.ErrValidation) {
		t.Errorf("expected ErrValidation, got: %v", err)
	}
}

func TestRecommendDevice_RejectsEmptySubject(t *testing.T) {
	ctx := context.Background()
	ai, students, devices, ramps, conversations, usage := newRecommendMocks("", nil)
	req := baseRecommendRequest
	req.Subject = ""

	_, err := inclusion.NewRecommendDevice(ai, students, devices, ramps, conversations, usage).Execute(ctx, req)

	if !errors.Is(err, providers.ErrValidation) {
		t.Errorf("expected ErrValidation, got: %v", err)
	}
}

func TestRecommendDevice_PersistsTurnWithMetadataWhenUserIDPresent(t *testing.T) {
	ctx := context.Background()
	aiResp := "Usá Timer Visual [DEVICE_ID:1]"
	ai, students, devices, ramps, conversations, usage := newRecommendMocks(aiResp, nil)
	var captured providers.AppendTurnParams
	conversations.AppendTurnFn = func(_ context.Context, p providers.AppendTurnParams) (int64, error) {
		captured = p
		return 99, nil
	}
	req := baseRecommendRequest
	req.UserID = 5

	got, err := inclusion.NewRecommendDevice(ai, students, devices, ramps, conversations, usage).Execute(ctx, req)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.ConversationID != 99 {
		t.Errorf("expected conversation_id 99, got %d", got.ConversationID)
	}
	if captured.Mode != "recommend" {
		t.Errorf("expected mode 'recommend', got %q", captured.Mode)
	}
	if captured.Metadata["recommended_device"] != int64(1) {
		t.Errorf("expected recommended_device 1 in metadata, got %v", captured.Metadata["recommended_device"])
	}
}
