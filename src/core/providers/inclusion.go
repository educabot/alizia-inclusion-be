package providers

import (
	"context"

	"github.com/google/uuid"

	"github.com/educabot/alizia-inclusion-be/src/core/entities"
)

type StudentProfileProvider interface {
	GetByStudentID(ctx context.Context, studentID int64) (*entities.StudentProfile, error)
	Upsert(ctx context.Context, profile *entities.StudentProfile) error
}

type StudentProvider interface {
	GetStudent(ctx context.Context, orgID uuid.UUID, id int64) (*entities.Student, error)
	ListByClassroom(ctx context.Context, orgID uuid.UUID, classroomID int64) ([]entities.Student, error)
}
