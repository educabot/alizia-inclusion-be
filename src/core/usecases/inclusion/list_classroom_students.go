package inclusion

import (
	"context"

	"github.com/google/uuid"

	"github.com/educabot/alizia-inclusion-be/src/core/entities"
	"github.com/educabot/alizia-inclusion-be/src/core/providers"
)

type ListClassroomStudentsRequest struct {
	OrgID       uuid.UUID
	ClassroomID int64
}

func (r ListClassroomStudentsRequest) Validate() error {
	if r.OrgID == uuid.Nil {
		return errOrgIDRequired
	}
	if r.ClassroomID <= 0 {
		return errClassroomIDRequired
	}
	return nil
}

type ListClassroomStudents interface {
	Execute(ctx context.Context, req ListClassroomStudentsRequest) ([]entities.Student, error)
}

type listClassroomStudentsImpl struct {
	students providers.StudentProvider
}

func NewListClassroomStudents(students providers.StudentProvider) ListClassroomStudents {
	return &listClassroomStudentsImpl{students: students}
}

func (uc *listClassroomStudentsImpl) Execute(ctx context.Context, req ListClassroomStudentsRequest) ([]entities.Student, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	return uc.students.ListByClassroom(ctx, req.OrgID, req.ClassroomID)
}
