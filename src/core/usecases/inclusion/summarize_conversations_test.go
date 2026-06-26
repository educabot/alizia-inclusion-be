package inclusion_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"gorm.io/datatypes"

	"github.com/educabot/alizia-inclusion-be/src/core/entities"
	"github.com/educabot/alizia-inclusion-be/src/core/providers"
	mocks "github.com/educabot/alizia-inclusion-be/src/core/providers/mocks"
	"github.com/educabot/alizia-inclusion-be/src/core/usecases/inclusion"
	"github.com/educabot/alizia-inclusion-be/src/testutil"
)

func studentPtr(id int64) *int64 { return &id }

func TestSummarizeConversations_SummarizesPending(t *testing.T) {
	conversations := new(mocks.MockConversationProvider)
	summaries := new(mocks.MockConversationSummaryProvider)
	ai := new(mocks.MockAIClient)
	ctx := context.Background()

	conv := entities.Conversation{
		ID:             7,
		OrganizationID: testutil.TestOrgID,
		UserID:         5,
		StudentID:      studentPtr(2),
		Messages: []entities.ConversationMessage{
			{Role: "user", Content: "en dictado se levanta y se va", Metadata: datatypes.JSON([]byte("{}"))},
			{Role: "assistant", Content: "probá fragmentar la consigna", Metadata: datatypes.JSON([]byte(`{"recommended_device":31}`))},
		},
	}
	conversations.On("ListPendingSummary", ctx, mock.Anything, mock.Anything).
		Return([]entities.Conversation{conv}, nil)

	ai.On("Chat", ctx, mock.Anything).Return(&providers.ChatResponse{
		Content: `{"summary":"Trabajamos la fuga en dictado: fragmentar la consigna.","topic_keywords":["dictado","atencion","fuga"]}`,
		Usage:   &providers.TokenUsage{TotalTokens: 120},
	}, nil)

	var capturedStudents, capturedDevices []int64
	var capturedSummary entities.ConversationSummary
	summaries.On("Upsert", ctx, mock.AnythingOfType("entities.ConversationSummary"), mock.Anything, mock.Anything).
		Run(func(args mock.Arguments) {
			capturedSummary = args.Get(1).(entities.ConversationSummary)
			capturedStudents = args.Get(2).([]int64)
			capturedDevices = args.Get(3).([]int64)
		}).Return(nil)

	uc := inclusion.NewSummarizeConversations(conversations, summaries, ai, nil, 20, 50)
	res, err := uc.Execute(ctx)

	require.NoError(t, err)
	assert.Equal(t, 1, res.Processed)
	assert.Equal(t, 0, res.Failed)
	assert.Equal(t, int64(7), capturedSummary.ConversationID)
	assert.Contains(t, capturedSummary.Summary, "dictado")
	assert.Equal(t, 120, capturedSummary.TokenCount)
	assert.ElementsMatch(t, []string{"dictado", "atencion", "fuga"}, []string(capturedSummary.TopicKeywords))
	assert.ElementsMatch(t, []int64{2}, capturedStudents)
	assert.ElementsMatch(t, []int64{31}, capturedDevices)
	summaries.AssertExpectations(t)
}

func TestSummarizeConversations_EmptyBatch(t *testing.T) {
	conversations := new(mocks.MockConversationProvider)
	summaries := new(mocks.MockConversationSummaryProvider)
	ai := new(mocks.MockAIClient)
	ctx := context.Background()

	conversations.On("ListPendingSummary", ctx, mock.Anything, mock.Anything).
		Return([]entities.Conversation{}, nil)

	uc := inclusion.NewSummarizeConversations(conversations, summaries, ai, nil, 20, 50)
	res, err := uc.Execute(ctx)

	require.NoError(t, err)
	assert.Equal(t, 0, res.Processed)
	assert.Equal(t, 0, res.Failed)
	ai.AssertNotCalled(t, "Chat", mock.Anything, mock.Anything)
	summaries.AssertNotCalled(t, "Upsert", mock.Anything, mock.Anything, mock.Anything, mock.Anything)
}

func TestSummarizeConversations_AIErrorSkipsAndContinues(t *testing.T) {
	conversations := new(mocks.MockConversationProvider)
	summaries := new(mocks.MockConversationSummaryProvider)
	ai := new(mocks.MockAIClient)
	ctx := context.Background()

	conv1 := entities.Conversation{ID: 1, OrganizationID: testutil.TestOrgID, UserID: 5,
		Messages: []entities.ConversationMessage{{Role: "user", Content: "uno", Metadata: datatypes.JSON([]byte("{}"))}}}
	conv2 := entities.Conversation{ID: 2, OrganizationID: testutil.TestOrgID, UserID: 5,
		Messages: []entities.ConversationMessage{{Role: "user", Content: "dos", Metadata: datatypes.JSON([]byte("{}"))}}}
	conversations.On("ListPendingSummary", ctx, mock.Anything, mock.Anything).
		Return([]entities.Conversation{conv1, conv2}, nil)

	// Primera llamada falla, segunda OK: el lote no se aborta.
	ai.On("Chat", ctx, mock.Anything).Return((*providers.ChatResponse)(nil), errors.New("boom")).Once()
	ai.On("Chat", ctx, mock.Anything).Return(&providers.ChatResponse{
		Content: `{"summary":"ok","topic_keywords":["x"]}`,
	}, nil).Once()

	summaries.On("Upsert", ctx, mock.Anything, mock.Anything, mock.Anything).Return(nil)

	uc := inclusion.NewSummarizeConversations(conversations, summaries, ai, nil, 20, 50)
	res, err := uc.Execute(ctx)

	require.NoError(t, err)
	assert.Equal(t, 1, res.Processed)
	assert.Equal(t, 1, res.Failed)
	summaries.AssertNumberOfCalls(t, "Upsert", 1)
}
