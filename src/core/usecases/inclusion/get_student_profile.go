package inclusion

import (
	"context"

	"github.com/google/uuid"

	"github.com/educabot/alizia-inclusion-be/src/core/entities"
	"github.com/educabot/alizia-inclusion-be/src/core/providers"
)

type GetStudentProfileRequest struct {
	OrgID     uuid.UUID
	StudentID int64
}

func (r GetStudentProfileRequest) Validate() error {
	if r.OrgID == uuid.Nil {
		return errOrgIDRequired
	}
	if r.StudentID <= 0 {
		return errStudentIDRequired
	}
	return nil
}

type GetStudentProfile interface {
	Execute(ctx context.Context, req GetStudentProfileRequest) (*entities.Student, error)
}

type getStudentProfileImpl struct {
	students providers.StudentProvider
}

func NewGetStudentProfile(students providers.StudentProvider) GetStudentProfile {
	return &getStudentProfileImpl{students: students}
}

func (uc *getStudentProfileImpl) Execute(ctx context.Context, req GetStudentProfileRequest) (*entities.Student, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	return uc.students.GetStudent(ctx, req.OrgID, req.StudentID)
}
