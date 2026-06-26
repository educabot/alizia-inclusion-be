package providers

import (
	"context"

	"github.com/google/uuid"

	"github.com/educabot/alizia-inclusion-be/src/core/entities"
)

type StudentNoteProvider interface {
	ListByStudent(ctx context.Context, orgID uuid.UUID, studentID int64) ([]entities.StudentNote, error)
	Create(ctx context.Context, note *entities.StudentNote) error
}
