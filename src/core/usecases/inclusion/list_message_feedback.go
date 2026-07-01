package inclusion

import (
	"context"

	"github.com/google/uuid"

	"github.com/educabot/alizia-inclusion-be/src/core/providers"
)

type ListMessageFeedbackRequest struct {
	OrgID uuid.UUID
	// Rating filtra por like/dislike; vacío devuelve todos.
	Rating string
}

func (r ListMessageFeedbackRequest) Validate() error {
	if r.OrgID == uuid.Nil {
		return errOrgIDRequired
	}
	if r.Rating != "" {
		if _, ok := validRatings[r.Rating]; !ok {
			return errInvalidRating
		}
	}
	return nil
}

type ListMessageFeedback interface {
	Execute(ctx context.Context, req ListMessageFeedbackRequest) ([]providers.MessageFeedbackReview, error)
}

type listMessageFeedbackImpl struct {
	feedback providers.MessageFeedbackProvider
}

func NewListMessageFeedback(feedback providers.MessageFeedbackProvider) ListMessageFeedback {
	return &listMessageFeedbackImpl{feedback: feedback}
}

func (uc *listMessageFeedbackImpl) Execute(ctx context.Context, req ListMessageFeedbackRequest) ([]providers.MessageFeedbackReview, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	return uc.feedback.List(ctx, req.OrgID, req.Rating)
}
