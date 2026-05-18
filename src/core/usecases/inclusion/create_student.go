package inclusion

import (
	"context"

	"github.com/google/uuid"

	"github.com/educabot/alizia-inclusion-be/src/core/entities"
	"github.com/educabot/alizia-inclusion-be/src/core/providers"
)

type CreateStudentRequest struct {
	OrgID       uuid.UUID
	ClassroomID int64
	Name        string
}

func (r CreateStudentRequest) Validate() error {
	if r.OrgID == uuid.Nil {
		return errOrgIDRequired
	}
	if r.ClassroomID <= 0 {
		return errClassroomIDRequired
	}
	if r.Name == "" {
		return errNameRequired
	}
	return nil
}

type CreateStudent interface {
	Execute(ctx context.Context, req CreateStudentRequest) (*entities.Student, error)
}

type createStudentImpl struct {
	students providers.StudentProvider
}

func NewCreateStudent(students providers.StudentProvider) CreateStudent {
	return &createStudentImpl{students: students}
}

func (uc *createStudentImpl) Execute(ctx context.Context, req CreateStudentRequest) (*entities.Student, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	student := &entities.Student{
		OrganizationID: req.OrgID,
		ClassroomID:    req.ClassroomID,
		Name:           req.Name,
	}

	if err := uc.students.Create(ctx, student); err != nil {
		return nil, err
	}
	return student, nil
}
