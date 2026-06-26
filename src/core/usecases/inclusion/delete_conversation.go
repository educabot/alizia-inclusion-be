package inclusion

import (
	"context"

	"github.com/google/uuid"

	"github.com/educabot/alizia-inclusion-be/src/core/providers"
)

type DeleteConversationRequest struct {
	OrgID          uuid.UUID
	ConversationID int64
}

func (r DeleteConversationRequest) Validate() error {
	if r.OrgID == uuid.Nil {
		return errOrgIDRequired
	}
	if r.ConversationID <= 0 {
		return errConversationIDRequired
	}
	return nil
}

type DeleteConversation interface {
	Execute(ctx context.Context, req DeleteConversationRequest) error
}

type deleteConversationImpl struct {
	conversations providers.ConversationProvider
}

func NewDeleteConversation(conversations providers.ConversationProvider) DeleteConversation {
	return &deleteConversationImpl{conversations: conversations}
}

func (uc *deleteConversationImpl) Execute(ctx context.Context, req DeleteConversationRequest) error {
	if err := req.Validate(); err != nil {
		return err
	}
	return uc.conversations.Delete(ctx, req.OrgID, req.ConversationID)
}
