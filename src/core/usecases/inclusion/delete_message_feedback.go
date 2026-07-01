package inclusion

import (
	"context"

	"github.com/google/uuid"

	"github.com/educabot/alizia-inclusion-be/src/core/providers"
)

type DeleteMessageFeedbackRequest struct {
	OrgID     uuid.UUID
	UserID    int64
	MessageID int64
}

func (r DeleteMessageFeedbackRequest) Validate() error {
	if r.OrgID == uuid.Nil {
		return errOrgIDRequired
	}
	if r.UserID <= 0 {
		return errUserIDRequired
	}
	if r.MessageID <= 0 {
		return errMessageIDRequired
	}
	return nil
}

type DeleteMessageFeedback interface {
	Execute(ctx context.Context, req DeleteMessageFeedbackRequest) error
}

type deleteMessageFeedbackImpl struct {
	feedback providers.MessageFeedbackProvider
}

func NewDeleteMessageFeedback(feedback providers.MessageFeedbackProvider) DeleteMessageFeedback {
	return &deleteMessageFeedbackImpl{feedback: feedback}
}

func (uc *deleteMessageFeedbackImpl) Execute(ctx context.Context, req DeleteMessageFeedbackRequest) error {
	if err := req.Validate(); err != nil {
		return err
	}
	return uc.feedback.Delete(ctx, req.OrgID, req.MessageID, req.UserID)
}
