package management

import (
	"context"

	"github.com/google/uuid"

	"github.com/educabot/alizia-inclusion-be/src/core/entities"
	"github.com/educabot/alizia-inclusion-be/src/core/providers"
)

type ListClassroomsRequest struct {
	OrgID uuid.UUID
}

func (r ListClassroomsRequest) Validate() error {
	if r.OrgID == uuid.Nil {
		return errOrgIDRequired
	}
	return nil
}

type ListClassrooms interface {
	Execute(ctx context.Context, req ListClassroomsRequest) ([]entities.Classroom, error)
}

type listClassroomsImpl struct {
	classrooms providers.ClassroomProvider
}

func NewListClassrooms(classrooms providers.ClassroomProvider) ListClassrooms {
	return &listClassroomsImpl{classrooms: classrooms}
}

func (uc *listClassroomsImpl) Execute(ctx context.Context, req ListClassroomsRequest) ([]entities.Classroom, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	return uc.classrooms.List(ctx, req.OrgID)
}
