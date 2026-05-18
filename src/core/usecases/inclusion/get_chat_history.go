package inclusion

import (
	"context"

	"github.com/google/uuid"

	"github.com/educabot/alizia-inclusion-be/src/core/entities"
	"github.com/educabot/alizia-inclusion-be/src/core/providers"
)

type GetChatHistoryRequest struct {
	OrgID  uuid.UUID
	UserID int64
	Mode   string
}

func (r GetChatHistoryRequest) Validate() error {
	if r.OrgID == uuid.Nil {
		return errOrgIDRequired
	}
	if r.UserID <= 0 {
		return errUserIDRequired
	}
	if r.Mode == "" {
		return errModeRequired
	}
	return nil
}

type GetChatHistory interface {
	Execute(ctx context.Context, req GetChatHistoryRequest) ([]entities.Conversation, error)
}

type getChatHistoryImpl struct {
	conversations providers.ConversationProvider
}

func NewGetChatHistory(conversations providers.ConversationProvider) GetChatHistory {
	return &getChatHistoryImpl{conversations: conversations}
}

func (uc *getChatHistoryImpl) Execute(ctx context.Context, req GetChatHistoryRequest) ([]entities.Conversation, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	return uc.conversations.ListByUser(ctx, req.OrgID, req.UserID, req.Mode)
}
