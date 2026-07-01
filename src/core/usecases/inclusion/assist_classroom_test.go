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

var assistClassroomBaseRequest = inclusion.AssistClassroomRequest{
	OrgID:       testutil.TestOrgID,
	ClassroomID: 1,
	Message:     "Lucas no puede concentrarse",
}

// assistClassroomMocks wires the providers a successful assist call exercises:
// devices.ListDevices + students.ListByClassroom always run, ai.Chat returns the
// supplied content (or error), AppendTurn / Record are left to each test to expect.
func assistClassroomMocks(t *testing.T, aiContent string, aiErr error) (
	*mockproviders.MockAIClient,
	*mockproviders.MockStudentProvider,
	*mockproviders.MockDeviceProvider,
	*mockproviders.MockConversationProvider,
	*mockproviders.MockAIUsageProvider,
) {
	t.Helper()
	ai := new(mockproviders.MockAIClient)
	students := new(mockproviders.MockStudentProvider)
	devices := new(mockproviders.MockDeviceProvider)
	conversations := new(mockproviders.MockConversationProvider)
	usage := new(mockproviders.MockAIUsageProvider)

	devices.On("ListDevices", mock.Anything, testutil.TestOrgID, (*int64)(nil)).
		Return([]entities.Device{testutil.NewDevice(1, 1, "Pictogramas")}, nil)
	students.On("List", mock.Anything, testutil.TestOrgID).
		Return([]entities.Student{testutil.NewStudent(1, 1, "Lucas")}, nil)
	if aiErr != nil {
		ai.On("Chat", mock.Anything, mock.AnythingOfType("[]providers.ChatMessage")).
			Return(nil, aiErr)
	} else {
		ai.On("Chat", mock.Anything, mock.AnythingOfType("[]providers.ChatMessage")).
			Return(&providers.ChatResponse{Content: aiContent}, nil)
	}
	return ai, students, devices, conversations, usage
}

func TestAssistClassroom_ReturnsAssistResponse(t *testing.T) {
	ai, students, devices, conversations, usage := assistClassroomMocks(t, "Podrias usar pictogramas [DEVICE_ID:1]", nil)

	got, err := inclusion.NewAssistClassroom(inclusion.AssistClassroomDeps{AI: ai, Students: students, Devices: devices, Conversations: conversations, Usage: usage}).
		Execute(context.Background(), assistClassroomBaseRequest)

	require.NoError(t, err)
	assert.NotEmpty(t, got.Response)
	ai.AssertExpectations(t)
}

func TestAssistClassroom_WorksInGuidedMode(t *testing.T) {
	ai, students, devices, conversations, usage := assistClassroomMocks(t, "Para que alumno necesitas la adaptacion?", nil)
	req := assistClassroomBaseRequest
	req.Mode = "guided"

	got, err := inclusion.NewAssistClassroom(inclusion.AssistClassroomDeps{AI: ai, Students: students, Devices: devices, Conversations: conversations, Usage: usage}).
		Execute(context.Background(), req)

	require.NoError(t, err)
	assert.NotEmpty(t, got.Response)
	ai.AssertExpectations(t)
}

func TestAssistClassroom_WrapsAIError(t *testing.T) {
	ai, students, devices, conversations, usage := assistClassroomMocks(t, "", errDB)

	_, err := inclusion.NewAssistClassroom(inclusion.AssistClassroomDeps{AI: ai, Students: students, Devices: devices, Conversations: conversations, Usage: usage}).
		Execute(context.Background(), assistClassroomBaseRequest)

	assert.ErrorIs(t, err, providers.ErrServiceUnavailable)
}

func TestAssistClassroom_RejectsNilOrgID(t *testing.T) {
	ai := new(mockproviders.MockAIClient)
	students := new(mockproviders.MockStudentProvider)
	devices := new(mockproviders.MockDeviceProvider)
	conversations := new(mockproviders.MockConversationProvider)
	usage := new(mockproviders.MockAIUsageProvider)
	req := assistClassroomBaseRequest
	req.OrgID = uuid.Nil

	_, err := inclusion.NewAssistClassroom(inclusion.AssistClassroomDeps{AI: ai, Students: students, Devices: devices, Conversations: conversations, Usage: usage}).
		Execute(context.Background(), req)

	assert.ErrorIs(t, err, providers.ErrValidation)
	ai.AssertNotCalled(t, "Chat", mock.Anything, mock.Anything)
	devices.AssertNotCalled(t, "ListDevices", mock.Anything, mock.Anything, mock.Anything)
}

func TestAssistClassroom_RejectsEmptyMessage(t *testing.T) {
	ai := new(mockproviders.MockAIClient)
	students := new(mockproviders.MockStudentProvider)
	devices := new(mockproviders.MockDeviceProvider)
	conversations := new(mockproviders.MockConversationProvider)
	usage := new(mockproviders.MockAIUsageProvider)
	req := assistClassroomBaseRequest
	req.Message = ""

	_, err := inclusion.NewAssistClassroom(inclusion.AssistClassroomDeps{AI: ai, Students: students, Devices: devices, Conversations: conversations, Usage: usage}).
		Execute(context.Background(), req)

	assert.ErrorIs(t, err, providers.ErrValidation)
	ai.AssertNotCalled(t, "Chat", mock.Anything, mock.Anything)
}

func TestAssistClassroom_PersistsTurnWhenUserIDPresent(t *testing.T) {
	ai, students, devices, conversations, usage := assistClassroomMocks(t, "Podrias usar pictogramas [DEVICE_ID:1]", nil)
	var captured providers.AppendTurnParams
	conversations.On("AppendTurn", mock.Anything, mock.AnythingOfType("providers.AppendTurnParams")).
		Run(func(args mock.Arguments) {
			p, ok := args.Get(1).(providers.AppendTurnParams)
			require.True(t, ok)
			captured = p
		}).
		Return(providers.AppendTurnResult{ConversationID: 42, AssistantMessageID: 420}, nil)
	req := assistClassroomBaseRequest
	req.UserID = 7

	got, err := inclusion.NewAssistClassroom(inclusion.AssistClassroomDeps{AI: ai, Students: students, Devices: devices, Conversations: conversations, Usage: usage}).
		Execute(context.Background(), req)

	require.NoError(t, err)
	assert.Equal(t, int64(42), got.ConversationID)
	assert.Equal(t, int64(420), got.MessageID)
	assert.Equal(t, int64(7), captured.UserID)
	assert.Equal(t, assistClassroomBaseRequest.Message, captured.UserContent)
	conversations.AssertExpectations(t)
}

func TestAssistClassroom_SkipsPersistenceWhenUserIDMissing(t *testing.T) {
	ai, students, devices, conversations, usage := assistClassroomMocks(t, "ok", nil)

	_, err := inclusion.NewAssistClassroom(inclusion.AssistClassroomDeps{AI: ai, Students: students, Devices: devices, Conversations: conversations, Usage: usage}).
		Execute(context.Background(), assistClassroomBaseRequest)

	require.NoError(t, err)
	conversations.AssertNotCalled(t, "AppendTurn", mock.Anything, mock.Anything)
}

func TestAssistClassroom_RecordsTokenUsageWhenPresent(t *testing.T) {
	ai := new(mockproviders.MockAIClient)
	students := new(mockproviders.MockStudentProvider)
	devices := new(mockproviders.MockDeviceProvider)
	conversations := new(mockproviders.MockConversationProvider)
	usage := new(mockproviders.MockAIUsageProvider)

	devices.On("ListDevices", mock.Anything, testutil.TestOrgID, (*int64)(nil)).
		Return([]entities.Device{testutil.NewDevice(1, 1, "Pictogramas")}, nil)
	students.On("List", mock.Anything, testutil.TestOrgID).
		Return([]entities.Student{testutil.NewStudent(1, 1, "Lucas")}, nil)
	ai.On("Chat", mock.Anything, mock.AnythingOfType("[]providers.ChatMessage")).
		Return(&providers.ChatResponse{
			Content: "ok",
			Usage:   &providers.TokenUsage{PromptTokens: 10, CompletionTokens: 5, TotalTokens: 15},
		}, nil)
	conversations.On("AppendTurn", mock.Anything, mock.AnythingOfType("providers.AppendTurnParams")).
		Return(providers.AppendTurnResult{ConversationID: 1, AssistantMessageID: 10}, nil)
	var captured providers.AIUsageRecord
	usage.On("Record", mock.Anything, mock.AnythingOfType("providers.AIUsageRecord")).
		Run(func(args mock.Arguments) {
			r, ok := args.Get(1).(providers.AIUsageRecord)
			require.True(t, ok)
			captured = r
		}).
		Return(nil)

	req := assistClassroomBaseRequest
	req.UserID = 7

	_, err := inclusion.NewAssistClassroom(inclusion.AssistClassroomDeps{AI: ai, Students: students, Devices: devices, Conversations: conversations, Usage: usage}).
		Execute(context.Background(), req)

	require.NoError(t, err)
	assert.Equal(t, 15, captured.TotalTokens)
	assert.Equal(t, "assist", captured.Mode)
	usage.AssertExpectations(t)
}

// TestAssistClassroom_AutoCreatesResourceFromAdaptationBlock: cuando el modelo emite un
// ADAPTATION_JSON sin id, el assist crea el recurso automáticamente con
// SourceConversationID seteado, y devuelve el id persistido en la respuesta.
func TestAssistClassroom_AutoCreatesResourceFromAdaptationBlock(t *testing.T) {
	content := `Probá fragmentar la consigna. [ADAPTATION_JSON:{"title":"Fragmentar consignas","type":"estrategia_aula","strategy":"dividir en pasos","steps":[{"orden":1,"texto":"Fragmentá la consigna"}]}]`
	ai, students, devices, conversations, usage := assistClassroomMocks(t, content, nil)
	conversations.On("AppendTurn", mock.Anything, mock.AnythingOfType("providers.AppendTurnParams")).
		Return(providers.AppendTurnResult{ConversationID: 42, AssistantMessageID: 420}, nil)

	adaptations := new(mockproviders.MockAdaptationProvider)
	var createdWith *entities.Adaptation
	adaptations.On("Create", mock.Anything, mock.AnythingOfType("*entities.Adaptation")).
		Run(func(args mock.Arguments) {
			a, ok := args.Get(1).(*entities.Adaptation)
			require.True(t, ok)
			a.ID = 100 // simula el id asignado por la BD
			createdWith = a
		}).Return(nil)
	adaptations.On("Get", mock.Anything, testutil.TestOrgID, int64(100)).
		Return(&entities.Adaptation{ID: 100, Title: "Fragmentar consignas", SourceConversationID: ptrInt64(42)}, nil)

	req := assistClassroomBaseRequest
	req.UserID = 7

	got, err := inclusion.NewAssistClassroom(inclusion.AssistClassroomDeps{
		AI: ai, Students: students, Devices: devices, Conversations: conversations, Usage: usage,
		Adaptations:      adaptations,
		CreateAdaptation: inclusion.NewCreateAdaptation(adaptations),
		UpdateAdaptation: inclusion.NewUpdateAdaptation(adaptations),
	}).Execute(context.Background(), req)

	require.NoError(t, err)
	require.NotNil(t, got.Adaptation)
	require.NotNil(t, got.Adaptation.ID)
	assert.Equal(t, int64(100), *got.Adaptation.ID)
	require.NotNil(t, createdWith)
	assert.Equal(t, int64(7), createdWith.TeacherID)
	require.NotNil(t, createdWith.SourceConversationID)
	assert.Equal(t, int64(42), *createdWith.SourceConversationID)
	adaptations.AssertExpectations(t)
}

// TestAssistClassroom_AutoUpdatesResourceWhenIDBelongsToConversation: si el modelo
// devuelve un id que pertenece a un recurso de ESTA conversación, se ACTUALIZA ese
// recurso en vez de crear uno nuevo.
func TestAssistClassroom_AutoUpdatesResourceWhenIDBelongsToConversation(t *testing.T) {
	content := `Ajusté el recurso. [ADAPTATION_JSON:{"id":100,"title":"Fragmentar consignas v2","type":"estrategia_aula","strategy":"dividir mejor","steps":[{"orden":1,"texto":"Fragmentá más fino"}]}]`
	ai, students, devices, conversations, usage := assistClassroomMocks(t, content, nil)
	conversations.On("AppendTurn", mock.Anything, mock.AnythingOfType("providers.AppendTurnParams")).
		Return(providers.AppendTurnResult{ConversationID: 55, AssistantMessageID: 550}, nil)

	adaptations := new(mockproviders.MockAdaptationProvider)
	// loadConversationResources: el recurso 100 ya existe en esta conversación.
	adaptations.On("List", mock.Anything, mock.AnythingOfType("providers.AdaptationFilter")).
		Return([]entities.Adaptation{{ID: 100, Title: "Fragmentar consignas", Status: "en_curso"}}, nil)
	adaptations.On("Get", mock.Anything, testutil.TestOrgID, int64(100)).
		Return(&entities.Adaptation{ID: 100, Title: "Fragmentar consignas"}, nil)
	adaptations.On("Update", mock.Anything, mock.AnythingOfType("*entities.Adaptation")).Return(nil)
	adaptations.On("SetDevices", mock.Anything, int64(100), mock.Anything).Return(nil)

	req := assistClassroomBaseRequest
	req.UserID = 7
	req.ConversationID = 55

	got, err := inclusion.NewAssistClassroom(inclusion.AssistClassroomDeps{
		AI: ai, Students: students, Devices: devices, Conversations: conversations, Usage: usage,
		Adaptations:      adaptations,
		CreateAdaptation: inclusion.NewCreateAdaptation(adaptations),
		UpdateAdaptation: inclusion.NewUpdateAdaptation(adaptations),
	}).Execute(context.Background(), req)

	require.NoError(t, err)
	require.NotNil(t, got.Adaptation)
	require.NotNil(t, got.Adaptation.ID)
	assert.Equal(t, int64(100), *got.Adaptation.ID)
	adaptations.AssertCalled(t, "Update", mock.Anything, mock.AnythingOfType("*entities.Adaptation"))
	adaptations.AssertNotCalled(t, "Create", mock.Anything, mock.Anything)
}
