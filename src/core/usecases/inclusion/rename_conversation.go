package inclusion

import (
	"context"

	"github.com/google/uuid"

	"github.com/educabot/alizia-inclusion-be/src/core/providers"
)

type RenameConversationRequest struct {
	OrgID          uuid.UUID
	ConversationID int64
	Title          string
}

func (r RenameConversationRequest) Validate() error {
	if r.OrgID == uuid.Nil {
		return errOrgIDRequired
	}
	if r.ConversationID <= 0 {
		return errConversationIDRequired
	}
	if r.Title == "" {
		return errTitleRequired
	}
	return nil
}

type RenameConversation interface {
	Execute(ctx context.Context, req RenameConversationRequest) error
}

type renameConversationImpl struct {
	conversations providers.ConversationProvider
}

func NewRenameConversation(conversations providers.ConversationProvider) RenameConversation {
	return &renameConversationImpl{conversations: conversations}
}

func (uc *renameConversationImpl) Execute(ctx context.Context, req RenameConversationRequest) error {
	if err := req.Validate(); err != nil {
		return err
	}
	return uc.conversations.Rename(ctx, req.OrgID, req.ConversationID, req.Title)
}
