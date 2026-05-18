package providers

import (
	"context"

	"github.com/google/uuid"

	"github.com/educabot/alizia-inclusion-be/src/core/entities"
)

type UserProvider interface {
	GetByID(ctx context.Context, orgID uuid.UUID, id int64) (*entities.User, error)
	GetByEmail(ctx context.Context, email string) (*entities.User, error)
	ListByRole(ctx context.Context, orgID uuid.UUID, role string) ([]entities.User, error)
}
