package inclusion

import (
	"context"

	"github.com/google/uuid"

	"github.com/educabot/alizia-inclusion-be/src/core/entities"
	"github.com/educabot/alizia-inclusion-be/src/core/providers"
)

type UpdateStudentRequest struct {
	OrgID       uuid.UUID
	StudentID   int64
	Name        *string
	ClassroomID *int64
}

func (r UpdateStudentRequest) Validate() error {
	if r.OrgID == uuid.Nil {
		return errOrgIDRequired
	}
	if r.StudentID <= 0 {
		return errStudentIDRequired
	}
	return nil
}

type UpdateStudent interface {
	Execute(ctx context.Context, req UpdateStudentRequest) (*entities.Student, error)
}

type updateStudentImpl struct {
	students providers.StudentProvider
}

func NewUpdateStudent(students providers.StudentProvider) UpdateStudent {
	return &updateStudentImpl{students: students}
}

func (uc *updateStudentImpl) Execute(ctx context.Context, req UpdateStudentRequest) (*entities.Student, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	existing, err := uc.students.GetStudent(ctx, req.OrgID, req.StudentID)
	if err != nil {
		return nil, err
	}

	if req.Name != nil {
		existing.Name = *req.Name
	}
	if req.ClassroomID != nil {
		existing.ClassroomID = req.ClassroomID
	}

	if err := uc.students.Update(ctx, existing); err != nil {
		return nil, err
	}
	return existing, nil
}
