package providers

import (
	"context"

	"github.com/google/uuid"

	"github.com/educabot/alizia-inclusion-be/src/core/entities"
)

type AdaptationProvider interface {
	List(ctx context.Context, orgID uuid.UUID, studentID *int64) ([]entities.Adaptation, error)
	Get(ctx context.Context, orgID uuid.UUID, id int64) (*entities.Adaptation, error)
	Create(ctx context.Context, adaptation *entities.Adaptation) error
	Update(ctx context.Context, adaptation *entities.Adaptation) error
	Delete(ctx context.Context, orgID uuid.UUID, id int64) error
}
