package providers

import (
	"context"

	"github.com/google/uuid"

	"github.com/educabot/alizia-inclusion-be/src/core/entities"
)

// ConversationSummaryProvider expone la lectura de resúmenes por entidad para la apertura
// (HU-1). La escritura (compactación al cerrar) se agrega en HU-5.
type ConversationSummaryProvider interface {
	RecentByStudent(ctx context.Context, orgID uuid.UUID, studentID int64, limit int) ([]entities.ConversationSummary, error)
	RecentByDevice(ctx context.Context, orgID uuid.UUID, deviceID int64, limit int) ([]entities.ConversationSummary, error)
	RecentByTopic(ctx context.Context, orgID uuid.UUID, keyword string, limit int) ([]entities.ConversationSummary, error)
}
