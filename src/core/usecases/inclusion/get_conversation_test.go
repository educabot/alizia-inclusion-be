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
	"github.com/educabot/alizia-inclusion-be/src/core/usecases/inclusion"
	mockproviders "github.com/educabot/alizia-inclusion-be/src/mocks/providers"
	"github.com/educabot/alizia-inclusion-be/src/testutil"
)

func TestGetConversation_ReturnsConversationWithMessages(t *testing.T) {
	// Arrange
	conversations := new(mockproviders.MockConversationProvider)
	ctx := context.Background()
	expected := &entities.Conversation{ID: 42, Mode: "assist", Messages: []entities.ConversationMessage{
		{Role: "user", Content: "Hola"},
		{Role: "assistant", Content: "Buenas"},
	}}
	conversations.On("GetWithMessages", ctx, testutil.TestOrgID, int64(42)).Return(expected, nil)

	// Act
	got, err := inclusion.NewGetConversation(conversations).Execute(ctx, inclusion.GetConversationRequest{
		OrgID:          testutil.TestOrgID,
		ConversationID: 42,
	})

	// Assert
	require.NoError(t, err)
	assert.Equal(t, int64(42), got.ID)
	assert.Len(t, got.Messages, 2)
	conversations.AssertExpectations(t)
}

func TestGetConversation_PropagatesNotFound(t *testing.T) {
	// Arrange
	conversations := new(mockproviders.MockConversationProvider)
	ctx := context.Background()
	conversations.On("GetWithMessages", ctx, testutil.TestOrgID, int64(999)).Return(nil, providers.ErrNotFound)

	// Act
	_, err := inclusion.NewGetConversation(conversations).Execute(ctx, inclusion.GetConversationRequest{
		OrgID:          testutil.TestOrgID,
		ConversationID: 999,
	})

	// Assert
	assert.ErrorIs(t, err, providers.ErrNotFound)
	conversations.AssertExpectations(t)
}

func TestGetConversation_RejectsNilOrgID(t *testing.T) {
	conversations := new(mockproviders.MockConversationProvider)

	_, err := inclusion.NewGetConversation(conversations).Execute(context.Background(), inclusion.GetConversationRequest{
		OrgID: uuid.Nil, ConversationID: 1,
	})

	assert.ErrorIs(t, err, providers.ErrValidation)
	conversations.AssertNotCalled(t, "GetWithMessages", mock.Anything, mock.Anything, mock.Anything)
}

func TestGetConversation_RejectsZeroConversationID(t *testing.T) {
	conversations := new(mockproviders.MockConversationProvider)

	_, err := inclusion.NewGetConversation(conversations).Execute(context.Background(), inclusion.GetConversationRequest{
		OrgID: testutil.TestOrgID, ConversationID: 0,
	})

	assert.ErrorIs(t, err, providers.ErrValidation)
	conversations.AssertNotCalled(t, "GetWithMessages", mock.Anything, mock.Anything, mock.Anything)
}
