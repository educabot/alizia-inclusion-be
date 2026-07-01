package providers

import (
	"context"
	"time"

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

// AppendTurnResult devuelve los ids persistidos del turno. AssistantMessageID es
// el id del mensaje del asistente, que el FE usa para anclar el feedback.
type AppendTurnResult struct {
	ConversationID     int64
	AssistantMessageID int64
}

type ConversationProvider interface {
	ListByUser(ctx context.Context, orgID uuid.UUID, userID int64, mode string) ([]entities.Conversation, error)
	AppendTurn(ctx context.Context, params AppendTurnParams) (AppendTurnResult, error)
	Delete(ctx context.Context, orgID uuid.UUID, id int64) error
	Rename(ctx context.Context, orgID uuid.UUID, id int64, title string) error
	// ListPendingSummary devuelve conversaciones inactivas (último mensaje <
	// idleBefore) sin resumen actualizado, con Messages precargados. Para el cron.
	ListPendingSummary(ctx context.Context, idleBefore time.Time, limit int) ([]entities.Conversation, error)
}
