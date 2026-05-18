package inclusion

import (
	"context"

	"github.com/google/uuid"

	"github.com/educabot/alizia-inclusion-be/src/core/entities"
	"github.com/educabot/alizia-inclusion-be/src/core/providers"
)

type ListStudentsRequest struct {
	OrgID       uuid.UUID
	ClassroomID *int64
}

func (r ListStudentsRequest) Validate() error {
	if r.OrgID == uuid.Nil {
		return errOrgIDRequired
	}
	return nil
}

type ListStudents interface {
	Execute(ctx context.Context, req ListStudentsRequest) ([]entities.Student, error)
}

type listStudentsImpl struct {
	students providers.StudentProvider
}

func NewListStudents(students providers.StudentProvider) ListStudents {
	return &listStudentsImpl{students: students}
}

func (uc *listStudentsImpl) Execute(ctx context.Context, req ListStudentsRequest) ([]entities.Student, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	if req.ClassroomID != nil {
		return uc.students.ListByClassroom(ctx, req.OrgID, *req.ClassroomID)
	}
	return uc.students.List(ctx, req.OrgID)
}
