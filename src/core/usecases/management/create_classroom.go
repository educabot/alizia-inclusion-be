package management

import (
	"context"

	"github.com/google/uuid"

	"github.com/educabot/alizia-inclusion-be/src/core/entities"
	"github.com/educabot/alizia-inclusion-be/src/core/providers"
)

type CreateClassroomRequest struct {
	OrgID   uuid.UUID
	Name    string
	Grade   *string
	Section *string
}

func (r CreateClassroomRequest) Validate() error {
	if r.OrgID == uuid.Nil {
		return errOrgIDRequired
	}
	if r.Name == "" {
		return errNameRequired
	}
	return nil
}

type CreateClassroom interface {
	Execute(ctx context.Context, req CreateClassroomRequest) (*entities.Classroom, error)
}

type createClassroomImpl struct {
	classrooms providers.ClassroomProvider
}

func NewCreateClassroom(classrooms providers.ClassroomProvider) CreateClassroom {
	return &createClassroomImpl{classrooms: classrooms}
}

func (uc *createClassroomImpl) Execute(ctx context.Context, req CreateClassroomRequest) (*entities.Classroom, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	classroom := &entities.Classroom{
		OrganizationID: req.OrgID,
		Name:           req.Name,
		Grade:          req.Grade,
		Section:        req.Section,
	}

	if err := uc.classrooms.Create(ctx, classroom); err != nil {
		return nil, err
	}
	return classroom, nil
}
