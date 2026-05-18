package management

import (
	"context"

	"github.com/google/uuid"

	"github.com/educabot/alizia-inclusion-be/src/core/entities"
	"github.com/educabot/alizia-inclusion-be/src/core/providers"
)

type UpdateClassroomRequest struct {
	OrgID       uuid.UUID
	ClassroomID int64
	Name        *string
	Grade       *string
	Section     *string
}

func (r UpdateClassroomRequest) Validate() error {
	if r.OrgID == uuid.Nil {
		return errOrgIDRequired
	}
	if r.ClassroomID <= 0 {
		return errClassroomIDRequired
	}
	return nil
}

type UpdateClassroom interface {
	Execute(ctx context.Context, req UpdateClassroomRequest) (*entities.Classroom, error)
}

type updateClassroomImpl struct {
	classrooms providers.ClassroomProvider
}

func NewUpdateClassroom(classrooms providers.ClassroomProvider) UpdateClassroom {
	return &updateClassroomImpl{classrooms: classrooms}
}

func (uc *updateClassroomImpl) Execute(ctx context.Context, req UpdateClassroomRequest) (*entities.Classroom, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	existing, err := uc.classrooms.Get(ctx, req.OrgID, req.ClassroomID)
	if err != nil {
		return nil, err
	}

	if req.Name != nil {
		existing.Name = *req.Name
	}
	if req.Grade != nil {
		existing.Grade = req.Grade
	}
	if req.Section != nil {
		existing.Section = req.Section
	}

	if err := uc.classrooms.Update(ctx, existing); err != nil {
		return nil, err
	}
	return existing, nil
}
