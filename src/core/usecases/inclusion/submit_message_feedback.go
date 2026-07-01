package inclusion

import (
	"context"

	"github.com/google/uuid"

	"github.com/educabot/alizia-inclusion-be/src/core/entities"
	"github.com/educabot/alizia-inclusion-be/src/core/providers"
)

var validRatings = map[string]struct{}{
	"like":    {},
	"dislike": {},
}

type SubmitMessageFeedbackRequest struct {
	OrgID     uuid.UUID
	UserID    int64
	MessageID int64
	Rating    string
	Comment   string
}

func (r SubmitMessageFeedbackRequest) Validate() error {
	if r.OrgID == uuid.Nil {
		return errOrgIDRequired
	}
	if r.UserID <= 0 {
		return errUserIDRequired
	}
	if r.MessageID <= 0 {
		return errMessageIDRequired
	}
	if _, ok := validRatings[r.Rating]; !ok {
		return errInvalidRating
	}
	return nil
}

type SubmitMessageFeedback interface {
	Execute(ctx context.Context, req SubmitMessageFeedbackRequest) (*entities.MessageFeedback, error)
}

type submitMessageFeedbackImpl struct {
	feedback providers.MessageFeedbackProvider
}

func NewSubmitMessageFeedback(feedback providers.MessageFeedbackProvider) SubmitMessageFeedback {
	return &submitMessageFeedbackImpl{feedback: feedback}
}

func (uc *submitMessageFeedbackImpl) Execute(ctx context.Context, req SubmitMessageFeedbackRequest) (*entities.MessageFeedback, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	// Derivamos el conversation_id desde el mensaje y validamos que pertenezca a la
	// org del token (evita feedback cruzado entre organizaciones).
	convID, msgOrg, err := uc.feedback.MessageContext(ctx, req.MessageID)
	if err != nil {
		return nil, err
	}
	if msgOrg != req.OrgID {
		return nil, providers.ErrNotFound
	}

	fb := &entities.MessageFeedback{
		ConversationMessageID: req.MessageID,
		ConversationID:        convID,
		OrganizationID:        req.OrgID,
		UserID:                req.UserID,
		Rating:                req.Rating,
		Comment:               req.Comment,
	}
	if err := uc.feedback.Upsert(ctx, fb); err != nil {
		return nil, err
	}
	return fb, nil
}
