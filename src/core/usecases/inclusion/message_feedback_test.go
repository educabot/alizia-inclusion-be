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

func TestSubmitMessageFeedback_DerivesConversationAndUpserts(t *testing.T) {
	feedback := new(mockproviders.MockMessageFeedbackProvider)
	ctx := context.Background()

	feedback.On("MessageContext", ctx, int64(7)).Return(int64(42), testutil.TestOrgID, nil)
	var captured *entities.MessageFeedback
	feedback.On("Upsert", ctx, mock.AnythingOfType("*entities.MessageFeedback")).
		Run(func(args mock.Arguments) {
			captured, _ = args.Get(1).(*entities.MessageFeedback)
		}).
		Return(nil)

	got, err := inclusion.NewSubmitMessageFeedback(feedback).Execute(ctx, inclusion.SubmitMessageFeedbackRequest{
		OrgID:     testutil.TestOrgID,
		UserID:    5,
		MessageID: 7,
		Rating:    "dislike",
		Comment:   "no aplica para este alumno",
	})

	require.NoError(t, err)
	require.NotNil(t, captured)
	assert.Equal(t, int64(42), captured.ConversationID)
	assert.Equal(t, int64(7), captured.ConversationMessageID)
	assert.Equal(t, "dislike", captured.Rating)
	assert.Equal(t, "no aplica para este alumno", captured.Comment)
	assert.Equal(t, testutil.TestOrgID, got.OrganizationID)
	feedback.AssertExpectations(t)
}

func TestSubmitMessageFeedback_RejectsInvalidRating(t *testing.T) {
	feedback := new(mockproviders.MockMessageFeedbackProvider)

	_, err := inclusion.NewSubmitMessageFeedback(feedback).Execute(context.Background(), inclusion.SubmitMessageFeedbackRequest{
		OrgID:     testutil.TestOrgID,
		UserID:    5,
		MessageID: 7,
		Rating:    "meh",
	})

	require.ErrorIs(t, err, providers.ErrValidation)
	feedback.AssertNotCalled(t, "MessageContext", mock.Anything, mock.Anything)
}

func TestSubmitMessageFeedback_RejectsCrossOrgMessage(t *testing.T) {
	feedback := new(mockproviders.MockMessageFeedbackProvider)
	ctx := context.Background()

	otherOrg := uuid.MustParse("11111111-1111-1111-1111-111111111111")
	feedback.On("MessageContext", ctx, int64(7)).Return(int64(42), otherOrg, nil)

	_, err := inclusion.NewSubmitMessageFeedback(feedback).Execute(ctx, inclusion.SubmitMessageFeedbackRequest{
		OrgID:     testutil.TestOrgID,
		UserID:    5,
		MessageID: 7,
		Rating:    "like",
	})

	require.ErrorIs(t, err, providers.ErrNotFound)
	feedback.AssertNotCalled(t, "Upsert", mock.Anything, mock.Anything)
}
