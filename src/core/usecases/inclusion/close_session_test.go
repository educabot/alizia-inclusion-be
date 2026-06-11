package inclusion_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"gorm.io/datatypes"

	"github.com/educabot/alizia-inclusion-be/src/core/entities"
	"github.com/educabot/alizia-inclusion-be/src/core/providers"
	mockproviders "github.com/educabot/alizia-inclusion-be/src/core/providers/mocks"
	"github.com/educabot/alizia-inclusion-be/src/core/usecases/inclusion"
	"github.com/educabot/alizia-inclusion-be/src/testutil"
)

type closeMocks struct {
	ai            *mockproviders.MockAIClient
	conversations *mockproviders.MockConversationProvider
	summaries     *mockproviders.MockConversationSummaryProvider
	usage         *mockproviders.MockAIUsageProvider
}

func newCloseMocks() closeMocks {
	m := closeMocks{
		ai:            new(mockproviders.MockAIClient),
		conversations: new(mockproviders.MockConversationProvider),
		summaries:     new(mockproviders.MockConversationSummaryProvider),
		usage:         new(mockproviders.MockAIUsageProvider),
	}
	// Per-turn usage trace (HU-6, T-6.5) is best-effort; optional in tests.
	m.usage.On("Record", mock.Anything, mock.AnythingOfType("providers.AIUsageRecord")).Return(nil).Maybe()
	return m
}

func (m closeMocks) usecase() inclusion.CloseSession {
	return inclusion.NewCloseSession(m.ai, m.conversations, m.summaries, m.usage)
}

func TestCloseSession_RejectsNilOrgID(t *testing.T) {
	m := newCloseMocks()

	_, err := m.usecase().Execute(context.Background(), inclusion.CloseSessionRequest{ConversationID: 1})

	assert.Error(t, err)
	m.conversations.AssertNotCalled(t, "GetWithMessages")
}

func TestCloseSession_RejectsZeroConversationID(t *testing.T) {
	m := newCloseMocks()

	_, err := m.usecase().Execute(context.Background(), inclusion.CloseSessionRequest{OrgID: testutil.TestOrgID})

	assert.Error(t, err)
	m.conversations.AssertNotCalled(t, "GetWithMessages")
}

func TestCloseSession_EmptyConversationSkipsLLMAndUpsert(t *testing.T) {
	// Arrange
	m := newCloseMocks()
	m.conversations.On("GetWithMessages", mock.Anything, testutil.TestOrgID, int64(42)).
		Return(&entities.Conversation{ID: 42}, nil)

	// Act
	got, err := m.usecase().Execute(context.Background(), inclusion.CloseSessionRequest{
		OrgID:          testutil.TestOrgID,
		UserID:         1,
		ConversationID: 42,
	})

	// Assert
	require.NoError(t, err)
	assert.Equal(t, int64(42), got.ConversationID)
	assert.Empty(t, got.Summary)
	m.ai.AssertNotCalled(t, "Chat", mock.Anything, mock.Anything)
	m.summaries.AssertNotCalled(t, "Upsert", mock.Anything, mock.Anything, mock.Anything, mock.Anything)
}

func TestCloseSession_CompactsAndUpsertsWithTags(t *testing.T) {
	// Arrange
	m := newCloseMocks()
	conv := &entities.Conversation{
		ID:        7,
		StudentID: ptrInt64(3),
		Messages: []entities.ConversationMessage{
			{Role: "user", Content: "¿Cómo ayudo a un alumno con TEA que se desregula?"},
			{
				Role:     "assistant",
				Content:  "Probemos pausas sensoriales.",
				Metadata: datatypes.JSON([]byte(`{"identified_student":3,"recommended_device":12}`)),
			},
		},
	}
	m.conversations.On("GetWithMessages", mock.Anything, testutil.TestOrgID, int64(7)).Return(conv, nil)
	m.ai.On("Chat", mock.Anything, mock.AnythingOfType("[]providers.ChatMessage")).
		Return(&providers.ChatResponse{
			Content: `Acá va: {"summary":"Hablamos de TEA y autorregulación. Próximo: pausas sensoriales.","topic_keywords":["TEA","Autorregulación","tea"]}`,
			Usage:   &providers.TokenUsage{TotalTokens: 120},
		}, nil)
	m.usage.On("Record", mock.Anything, mock.Anything).Return(nil)

	var captured *entities.ConversationSummary
	var capturedStudents, capturedDevices []int64
	m.summaries.On("Upsert", mock.Anything, mock.AnythingOfType("*entities.ConversationSummary"), mock.Anything, mock.Anything).
		Run(func(args mock.Arguments) {
			captured, _ = args.Get(1).(*entities.ConversationSummary)
			capturedStudents, _ = args.Get(2).([]int64)
			capturedDevices, _ = args.Get(3).([]int64)
		}).Return(nil)

	// Act
	got, err := m.usecase().Execute(context.Background(), inclusion.CloseSessionRequest{
		OrgID:          testutil.TestOrgID,
		UserID:         1,
		ConversationID: 7,
	})

	// Assert
	require.NoError(t, err)
	assert.Equal(t, int64(7), got.ConversationID)
	assert.Contains(t, got.Summary, "pausas sensoriales")
	// Keywords are lowercased and deduplicated ("TEA" and "tea" collapse into one).
	assert.Equal(t, []string{"tea", "autorregulación"}, got.TopicKeywords)
	assert.Equal(t, []int64{3}, got.StudentIDs)
	assert.Equal(t, []int64{12}, got.DeviceIDs)

	require.NotNil(t, captured)
	assert.Equal(t, int64(7), captured.ConversationID)
	assert.Equal(t, []int64{3}, capturedStudents)
	assert.Equal(t, []int64{12}, capturedDevices)
	m.conversations.AssertExpectations(t)
	m.ai.AssertExpectations(t)
	m.summaries.AssertExpectations(t)
}

func TestCloseSession_FallsBackToRawContentWhenNotJSON(t *testing.T) {
	// Arrange
	m := newCloseMocks()
	conv := &entities.Conversation{
		ID:       9,
		Messages: []entities.ConversationMessage{{Role: "user", Content: "Hola"}},
	}
	m.conversations.On("GetWithMessages", mock.Anything, testutil.TestOrgID, int64(9)).Return(conv, nil)
	m.ai.On("Chat", mock.Anything, mock.AnythingOfType("[]providers.ChatMessage")).
		Return(&providers.ChatResponse{Content: "Resumen en texto plano sin JSON."}, nil)
	m.summaries.On("Upsert", mock.Anything, mock.AnythingOfType("*entities.ConversationSummary"), mock.Anything, mock.Anything).
		Return(nil)

	// Act
	got, err := m.usecase().Execute(context.Background(), inclusion.CloseSessionRequest{
		OrgID:          testutil.TestOrgID,
		UserID:         1,
		ConversationID: 9,
	})

	// Assert
	require.NoError(t, err)
	assert.Equal(t, "Resumen en texto plano sin JSON.", got.Summary)
	assert.Empty(t, got.TopicKeywords)
	m.summaries.AssertExpectations(t)
}

func TestCloseSession_PropagatesConversationNotFound(t *testing.T) {
	// Arrange
	m := newCloseMocks()
	m.conversations.On("GetWithMessages", mock.Anything, testutil.TestOrgID, int64(404)).
		Return(nil, providers.ErrNotFound)

	// Act
	_, err := m.usecase().Execute(context.Background(), inclusion.CloseSessionRequest{
		OrgID:          testutil.TestOrgID,
		UserID:         1,
		ConversationID: 404,
	})

	// Assert
	assert.ErrorIs(t, err, providers.ErrNotFound)
	m.ai.AssertNotCalled(t, "Chat", mock.Anything, mock.Anything)
}
