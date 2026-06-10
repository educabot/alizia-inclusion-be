package providers

import (
	"context"

	"github.com/google/uuid"

	"github.com/educabot/alizia-inclusion-be/src/core/entities"
)

// ConversationSummaryProvider expone la lectura de resúmenes por entidad para la apertura
// (HU-1) y la escritura por compactación al cerrar (HU-5).
type ConversationSummaryProvider interface {
	RecentByStudent(ctx context.Context, orgID uuid.UUID, studentID int64, limit int) ([]entities.ConversationSummary, error)
	RecentByDevice(ctx context.Context, orgID uuid.UUID, deviceID int64, limit int) ([]entities.ConversationSummary, error)
	RecentByTopic(ctx context.Context, orgID uuid.UUID, keyword string, limit int) ([]entities.ConversationSummary, error)
	// Upsert guarda/actualiza el resumen compactado y revincula el resumen a sus
	// entidades (alumnos / devices). Idempotente por conversation_id.
	Upsert(ctx context.Context, summary *entities.ConversationSummary, studentIDs, deviceIDs []int64) error
}
