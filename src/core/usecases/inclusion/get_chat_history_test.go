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
	mockproviders "github.com/educabot/alizia-inclusion-be/src/mocks/providers"
	"github.com/educabot/alizia-inclusion-be/src/core/usecases/inclusion"
	"github.com/educabot/alizia-inclusion-be/src/testutil"
)

func TestGetChatHistory_ReturnsConversations(t *testing.T) {
	conversations := new(mockproviders.MockConversationProvider)
	ctx := context.Background()
	expected := []entities.Conversation{
		testutil.NewConversation(1, 1, "recommend"),
		testutil.NewConversation(2, 1, "recommend"),
	}
	conversations.On("ListByUser", ctx, testutil.TestOrgID, int64(1), "recommend").Return(expected, nil)

	got, err := inclusion.NewGetChatHistory(conversations).Execute(ctx, inclusion.GetChatHistoryRequest{
		OrgID:  testutil.TestOrgID,
		UserID: 1,
		Mode:   "recommend",
	})

	require.NoError(t, err)
	assert.Len(t, got, 2)
	conversations.AssertExpectations(t)
}

func TestGetChatHistory_RejectsNilOrgID(t *testing.T) {
	conversations := new(mockproviders.MockConversationProvider)

	_, err := inclusion.NewGetChatHistory(conversations).Execute(context.Background(), inclusion.GetChatHistoryRequest{
		OrgID: uuid.Nil, UserID: 1, Mode: "recommend",
	})

	assert.ErrorIs(t, err, providers.ErrValidation)
	conversations.AssertNotCalled(t, "ListByUser", mock.Anything, mock.Anything, mock.Anything, mock.Anything)
}

func TestGetChatHistory_RejectsZeroUserID(t *testing.T) {
	conversations := new(mockproviders.MockConversationProvider)

	_, err := inclusion.NewGetChatHistory(conversations).Execute(context.Background(), inclusion.GetChatHistoryRequest{
		OrgID: testutil.TestOrgID, UserID: 0, Mode: "recommend",
	})

	assert.ErrorIs(t, err, providers.ErrValidation)
	conversations.AssertNotCalled(t, "ListByUser", mock.Anything, mock.Anything, mock.Anything, mock.Anything)
}

func TestGetChatHistory_RejectsEmptyMode(t *testing.T) {
	conversations := new(mockproviders.MockConversationProvider)

	_, err := inclusion.NewGetChatHistory(conversations).Execute(context.Background(), inclusion.GetChatHistoryRequest{
		OrgID: testutil.TestOrgID, UserID: 1, Mode: "",
	})

	assert.ErrorIs(t, err, providers.ErrValidation)
	conversations.AssertNotCalled(t, "ListByUser", mock.Anything, mock.Anything, mock.Anything, mock.Anything)
}
