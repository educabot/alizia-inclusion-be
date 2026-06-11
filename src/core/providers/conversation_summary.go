package providers

import (
	"context"

	"github.com/google/uuid"

	"github.com/educabot/alizia-inclusion-be/src/core/entities"
)

// ConversationSummaryProvider exposes summary reads by entity on conversation open (HU-1)
// and compacted summary writes on close (HU-5).
type ConversationSummaryProvider interface {
	RecentByStudent(ctx context.Context, orgID uuid.UUID, studentID int64, limit int) ([]entities.ConversationSummary, error)
	RecentByDevice(ctx context.Context, orgID uuid.UUID, deviceID int64, limit int) ([]entities.ConversationSummary, error)
	RecentByTopic(ctx context.Context, orgID uuid.UUID, keyword string, limit int) ([]entities.ConversationSummary, error)
	// Upsert saves or updates the compacted summary and re-links it to its
	// related entities (students / devices). Idempotent by conversation_id.
	Upsert(ctx context.Context, summary *entities.ConversationSummary, studentIDs, deviceIDs []int64) error
}
