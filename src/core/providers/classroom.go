package providers

import (
	"context"

	"github.com/google/uuid"

	"github.com/educabot/alizia-inclusion-be/src/core/entities"
)

type ClassroomProvider interface {
	List(ctx context.Context, orgID uuid.UUID) ([]entities.Classroom, error)
	Get(ctx context.Context, orgID uuid.UUID, id int64) (*entities.Classroom, error)
	Create(ctx context.Context, classroom *entities.Classroom) error
	Update(ctx context.Context, classroom *entities.Classroom) error
	Delete(ctx context.Context, orgID uuid.UUID, id int64) error
}
