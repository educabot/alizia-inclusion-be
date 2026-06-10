package providers

import (
	"context"

	"github.com/google/uuid"

	"github.com/educabot/alizia-inclusion-be/src/core/entities"
)

type AppendTurnParams struct {
	ConversationID   int64
	OrgID            uuid.UUID
	UserID           int64
	Mode             string
	StudentID        *int64
	UserContent      string
	AssistantContent string
	Metadata         map[string]any
}

type ConversationProvider interface {
	ListByUser(ctx context.Context, orgID uuid.UUID, userID int64, mode string) ([]entities.Conversation, error)
	AppendTurn(ctx context.Context, params AppendTurnParams) (int64, error)
	// GetWithMessages trae una conversación con sus mensajes ordenados, para
	// compactarla al cerrar (HU-5).
	GetWithMessages(ctx context.Context, orgID uuid.UUID, conversationID int64) (*entities.Conversation, error)
}
