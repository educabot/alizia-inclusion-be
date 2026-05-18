package inclusion

import (
	"context"

	"github.com/google/uuid"

	"github.com/educabot/alizia-inclusion-be/src/core/providers"
)

type DeleteStudentRequest struct {
	OrgID     uuid.UUID
	StudentID int64
}

func (r DeleteStudentRequest) Validate() error {
	if r.OrgID == uuid.Nil {
		return errOrgIDRequired
	}
	if r.StudentID <= 0 {
		return errStudentIDRequired
	}
	return nil
}

type DeleteStudent interface {
	Execute(ctx context.Context, req DeleteStudentRequest) error
}

type deleteStudentImpl struct {
	students providers.StudentProvider
}

func NewDeleteStudent(students providers.StudentProvider) DeleteStudent {
	return &deleteStudentImpl{students: students}
}

func (uc *deleteStudentImpl) Execute(ctx context.Context, req DeleteStudentRequest) error {
	if err := req.Validate(); err != nil {
		return err
	}
	return uc.students.Delete(ctx, req.OrgID, req.StudentID)
}
