package inclusion

import (
	"context"

	"github.com/google/uuid"

	"github.com/educabot/alizia-inclusion-be/src/core/entities"
	"github.com/educabot/alizia-inclusion-be/src/core/providers"
)

type GetConversationRequest struct {
	OrgID          uuid.UUID
	ConversationID int64
}

func (r GetConversationRequest) Validate() error {
	if r.OrgID == uuid.Nil {
		return errOrgIDRequired
	}
	if r.ConversationID <= 0 {
		return errConversationIDRequired
	}
	return nil
}

// GetConversation loads a single conversation with its messages, scoped to the org.
// Used by the frontend to resume the conversation that originated a saved resource
// (adaptation.source_conversation_id).
type GetConversation interface {
	Execute(ctx context.Context, req GetConversationRequest) (*entities.Conversation, error)
}

type getConversationImpl struct {
	conversations providers.ConversationProvider
}

func NewGetConversation(conversations providers.ConversationProvider) GetConversation {
	return &getConversationImpl{conversations: conversations}
}

func (uc *getConversationImpl) Execute(ctx context.Context, req GetConversationRequest) (*entities.Conversation, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	return uc.conversations.GetWithMessages(ctx, req.OrgID, req.ConversationID)
}
