package inclusion_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/educabot/alizia-inclusion-be/src/core/entities"
	"github.com/educabot/alizia-inclusion-be/src/core/providers"
	mockproviders "github.com/educabot/alizia-inclusion-be/src/core/providers/mocks"
	"github.com/educabot/alizia-inclusion-be/src/core/usecases/inclusion"
	"github.com/educabot/alizia-inclusion-be/src/testutil"
)

var baseRecommendRequest = inclusion.RecommendDeviceRequest{
	OrgID:     testutil.TestOrgID,
	StudentID: 1,
	Subject:   "Matematicas",
	Objective: "Sumar fracciones",
}

// recommendDeviceMocks wires the providers a successful recommend call exercises.
func recommendDeviceMocks(t *testing.T, aiContent string, aiErr error) (
	*mockproviders.MockAIClient,
	*mockproviders.MockStudentProvider,
	*mockproviders.MockDeviceProvider,
	*mockproviders.MockRampProvider,
	*mockproviders.MockConversationProvider,
	*mockproviders.MockAIUsageProvider,
) {
	t.Helper()
	ai := new(mockproviders.MockAIClient)
	students := new(mockproviders.MockStudentProvider)
	devices := new(mockproviders.MockDeviceProvider)
	ramps := new(mockproviders.MockRampProvider)
	conversations := new(mockproviders.MockConversationProvider)
	usage := new(mockproviders.MockAIUsageProvider)

	student := testutil.NewStudentWithProfile(1, 1, "Lucas", []string{"distraccion"})
	device := testutil.NewDevice(1, 1, "Timer Visual")

	students.On("GetStudent", mock.Anything, testutil.TestOrgID, int64(1)).Return(&student, nil)
	devices.On("ListDevices", mock.Anything, testutil.TestOrgID, (*int64)(nil)).
		Return([]entities.Device{device}, nil)
	if aiErr != nil {
		ai.On("Chat", mock.Anything, mock.AnythingOfType("[]providers.ChatMessage")).
			Return(nil, aiErr)
	} else {
		ai.On("Chat", mock.Anything, mock.AnythingOfType("[]providers.ChatMessage")).
			Return(&providers.ChatResponse{Content: aiContent}, nil)
	}
	return ai, students, devices, ramps, conversations, usage
}

func TestRecommendDevice_ReturnsRecommendationWithDevice(t *testing.T) {
	aiResp := `Recomiendo el Timer Visual [DEVICE_ID:1] para ayudar con la distraccion.
[ADAPTATION_JSON:{"title":"Timer para fracciones","type":"actividad_adaptada","strategy":"Usar timer","device_ids":[1],"device_names":["Timer Visual"]}]`
	ai, students, devices, ramps, conversations, usage := recommendDeviceMocks(t, aiResp, nil)

	got, err := inclusion.NewRecommendDevice(ai, students, devices, ramps, conversations, usage).
		Execute(context.Background(), baseRecommendRequest)

	require.NoError(t, err)
	assert.NotEmpty(t, got.Response)
	require.NotNil(t, got.DeviceID)
	assert.Equal(t, int64(1), *got.DeviceID)
	require.NotNil(t, got.Adaptation)
	assert.Equal(t, "Timer para fracciones", got.Adaptation.Title)
	ai.AssertExpectations(t)
	students.AssertExpectations(t)
}

func TestRecommendDevice_HandlesResponseWithoutMarkers(t *testing.T) {
	ai, students, devices, ramps, conversations, usage := recommendDeviceMocks(t, "Respuesta sin marcadores", nil)

	got, err := inclusion.NewRecommendDevice(ai, students, devices, ramps, conversations, usage).
		Execute(context.Background(), baseRecommendRequest)

	require.NoError(t, err)
	assert.Nil(t, got.DeviceID)
	assert.Nil(t, got.Adaptation)
}

func TestRecommendDevice_WrapsAIErrorAsServiceUnavailable(t *testing.T) {
	ai, students, devices, ramps, conversations, usage := recommendDeviceMocks(t, "", errDB)

	_, err := inclusion.NewRecommendDevice(ai, students, devices, ramps, conversations, usage).
		Execute(context.Background(), baseRecommendRequest)

	assert.ErrorIs(t, err, providers.ErrServiceUnavailable)
}

func TestRecommendDevice_RejectsNilOrgID(t *testing.T) {
	ai := new(mockproviders.MockAIClient)
	students := new(mockproviders.MockStudentProvider)
	devices := new(mockproviders.MockDeviceProvider)
	ramps := new(mockproviders.MockRampProvider)
	conversations := new(mockproviders.MockConversationProvider)
	usage := new(mockproviders.MockAIUsageProvider)
	req := baseRecommendRequest
	req.OrgID = uuid.Nil

	_, err := inclusion.NewRecommendDevice(ai, students, devices, ramps, conversations, usage).
		Execute(context.Background(), req)

	assert.ErrorIs(t, err, providers.ErrValidation)
	ai.AssertNotCalled(t, "Chat", mock.Anything, mock.Anything)
	students.AssertNotCalled(t, "GetStudent", mock.Anything, mock.Anything, mock.Anything)
}

func TestRecommendDevice_RejectsZeroStudentID(t *testing.T) {
	ai := new(mockproviders.MockAIClient)
	students := new(mockproviders.MockStudentProvider)
	devices := new(mockproviders.MockDeviceProvider)
	ramps := new(mockproviders.MockRampProvider)
	conversations := new(mockproviders.MockConversationProvider)
	usage := new(mockproviders.MockAIUsageProvider)
	req := baseRecommendRequest
	req.StudentID = 0

	_, err := inclusion.NewRecommendDevice(ai, students, devices, ramps, conversations, usage).
		Execute(context.Background(), req)

	assert.ErrorIs(t, err, providers.ErrValidation)
	ai.AssertNotCalled(t, "Chat", mock.Anything, mock.Anything)
}

func TestRecommendDevice_RejectsEmptySubject(t *testing.T) {
	ai := new(mockproviders.MockAIClient)
	students := new(mockproviders.MockStudentProvider)
	devices := new(mockproviders.MockDeviceProvider)
	ramps := new(mockproviders.MockRampProvider)
	conversations := new(mockproviders.MockConversationProvider)
	usage := new(mockproviders.MockAIUsageProvider)
	req := baseRecommendRequest
	req.Subject = ""

	_, err := inclusion.NewRecommendDevice(ai, students, devices, ramps, conversations, usage).
		Execute(context.Background(), req)

	assert.ErrorIs(t, err, providers.ErrValidation)
	ai.AssertNotCalled(t, "Chat", mock.Anything, mock.Anything)
}

func TestRecommendDevice_PersistsTurnWithMetadataWhenUserIDPresent(t *testing.T) {
	ai, students, devices, ramps, conversations, usage := recommendDeviceMocks(t, "Usá Timer Visual [DEVICE_ID:1]", nil)
	var captured providers.AppendTurnParams
	conversations.On("AppendTurn", mock.Anything, mock.AnythingOfType("providers.AppendTurnParams")).
		Run(func(args mock.Arguments) {
			captured = args.Get(1).(providers.AppendTurnParams)
		}).
		Return(int64(99), nil)
	req := baseRecommendRequest
	req.UserID = 5

	got, err := inclusion.NewRecommendDevice(ai, students, devices, ramps, conversations, usage).
		Execute(context.Background(), req)

	require.NoError(t, err)
	assert.Equal(t, int64(99), got.ConversationID)
	assert.Equal(t, "recommend", captured.Mode)
	assert.Equal(t, int64(1), captured.Metadata["recommended_device"])
	conversations.AssertExpectations(t)
}
