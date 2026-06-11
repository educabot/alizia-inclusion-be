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

type recommendMocks struct {
	ai            *mockproviders.MockAIClient
	students      *mockproviders.MockStudentProvider
	devices       *mockproviders.MockDeviceProvider
	ramps         *mockproviders.MockRampProvider
	conversations *mockproviders.MockConversationProvider
	usage         *mockproviders.MockAIUsageProvider
}

func newRecommendMocks() recommendMocks {
	m := recommendMocks{
		ai:            new(mockproviders.MockAIClient),
		students:      new(mockproviders.MockStudentProvider),
		devices:       new(mockproviders.MockDeviceProvider),
		ramps:         new(mockproviders.MockRampProvider),
		conversations: new(mockproviders.MockConversationProvider),
		usage:         new(mockproviders.MockAIUsageProvider),
	}
	// La traza por turno (HU-6, T-6.5) se graba best-effort; opcional para los tests.
	m.usage.On("Record", mock.Anything, mock.AnythingOfType("providers.AIUsageRecord")).Return(nil).Maybe()
	return m
}

func (m recommendMocks) usecase() inclusion.RecommendDevice {
	return inclusion.NewRecommendDevice(m.ai, m.students, m.devices, m.ramps, m.conversations, m.usage)
}

// recommendDeviceMocks wires the providers a successful recommend call exercises.
func recommendDeviceMocks(t *testing.T, aiContent string, aiErr error) recommendMocks {
	t.Helper()
	m := newRecommendMocks()
	student := testutil.NewStudentWithProfile(1, 1, "Lucas", []string{"distraccion"})
	device := testutil.NewDevice(1, 1, "Timer Visual")

	m.students.On("GetStudent", mock.Anything, testutil.TestOrgID, int64(1)).Return(&student, nil)
	m.devices.On("ListDevices", mock.Anything, testutil.TestOrgID, (*int64)(nil)).
		Return([]entities.Device{device}, nil)
	if aiErr != nil {
		m.ai.On("Chat", mock.Anything, mock.AnythingOfType("[]providers.ChatMessage")).
			Return(nil, aiErr)
	} else {
		m.ai.On("Chat", mock.Anything, mock.AnythingOfType("[]providers.ChatMessage")).
			Return(&providers.ChatResponse{Content: aiContent}, nil)
	}
	return m
}

func TestRecommendDevice_ReturnsRecommendationWithDevice(t *testing.T) {
	aiResp := `Recomiendo el Timer Visual [DEVICE_ID:1] para ayudar con la distraccion.
[ADAPTATION_JSON:{"title":"Timer para fracciones","type":"actividad_adaptada","strategy":"Usar timer","device_ids":[1],"device_names":["Timer Visual"]}]`
	m := recommendDeviceMocks(t, aiResp, nil)

	got, err := m.usecase().Execute(context.Background(), baseRecommendRequest)

	require.NoError(t, err)
	assert.NotEmpty(t, got.Response)
	require.NotNil(t, got.DeviceID)
	assert.Equal(t, int64(1), *got.DeviceID)
	require.NotNil(t, got.Adaptation)
	assert.Equal(t, "Timer para fracciones", got.Adaptation.Title)
	m.ai.AssertExpectations(t)
	m.students.AssertExpectations(t)
}

func TestRecommendDevice_HandlesResponseWithoutMarkers(t *testing.T) {
	m := recommendDeviceMocks(t, "Respuesta sin marcadores", nil)

	got, err := m.usecase().Execute(context.Background(), baseRecommendRequest)

	require.NoError(t, err)
	assert.Nil(t, got.DeviceID)
	assert.Nil(t, got.Adaptation)
}

func TestRecommendDevice_WrapsAIErrorAsServiceUnavailable(t *testing.T) {
	m := recommendDeviceMocks(t, "", errDB)

	_, err := m.usecase().Execute(context.Background(), baseRecommendRequest)

	assert.ErrorIs(t, err, providers.ErrServiceUnavailable)
}

func TestRecommendDevice_RejectsNilOrgID(t *testing.T) {
	m := newRecommendMocks()
	req := baseRecommendRequest
	req.OrgID = uuid.Nil

	_, err := m.usecase().Execute(context.Background(), req)

	assert.ErrorIs(t, err, providers.ErrValidation)
	m.ai.AssertNotCalled(t, "Chat", mock.Anything, mock.Anything)
	m.students.AssertNotCalled(t, "GetStudent", mock.Anything, mock.Anything, mock.Anything)
}

func TestRecommendDevice_RejectsZeroStudentID(t *testing.T) {
	m := newRecommendMocks()
	req := baseRecommendRequest
	req.StudentID = 0

	_, err := m.usecase().Execute(context.Background(), req)

	assert.ErrorIs(t, err, providers.ErrValidation)
	m.ai.AssertNotCalled(t, "Chat", mock.Anything, mock.Anything)
}

func TestRecommendDevice_RejectsEmptySubject(t *testing.T) {
	m := newRecommendMocks()
	req := baseRecommendRequest
	req.Subject = ""

	_, err := m.usecase().Execute(context.Background(), req)

	assert.ErrorIs(t, err, providers.ErrValidation)
	m.ai.AssertNotCalled(t, "Chat", mock.Anything, mock.Anything)
}

func TestRecommendDevice_PersistsTurnWithMetadataWhenUserIDPresent(t *testing.T) {
	m := recommendDeviceMocks(t, "Usá Timer Visual [DEVICE_ID:1]", nil)
	var captured providers.AppendTurnParams
	m.conversations.On("AppendTurn", mock.Anything, mock.AnythingOfType("providers.AppendTurnParams")).
		Run(func(args mock.Arguments) {
			p, ok := args.Get(1).(providers.AppendTurnParams)
			require.True(t, ok)
			captured = p
		}).
		Return(int64(99), nil)
	req := baseRecommendRequest
	req.UserID = 5

	got, err := m.usecase().Execute(context.Background(), req)

	require.NoError(t, err)
	assert.Equal(t, int64(99), got.ConversationID)
	assert.Equal(t, "recommend", captured.Mode)
	assert.Equal(t, int64(1), captured.Metadata["recommended_device"])
	m.conversations.AssertExpectations(t)
}
