package management

import (
	"context"

	"github.com/google/uuid"

	"github.com/educabot/alizia-inclusion-be/src/core/providers"
)

type DeleteClassroomRequest struct {
	OrgID       uuid.UUID
	ClassroomID int64
}

func (r DeleteClassroomRequest) Validate() error {
	if r.OrgID == uuid.Nil {
		return errOrgIDRequired
	}
	if r.ClassroomID <= 0 {
		return errClassroomIDRequired
	}
	return nil
}

type DeleteClassroom interface {
	Execute(ctx context.Context, req DeleteClassroomRequest) error
}

type deleteClassroomImpl struct {
	classrooms providers.ClassroomProvider
}

func NewDeleteClassroom(classrooms providers.ClassroomProvider) DeleteClassroom {
	return &deleteClassroomImpl{classrooms: classrooms}
}

func (uc *deleteClassroomImpl) Execute(ctx context.Context, req DeleteClassroomRequest) error {
	if err := req.Validate(); err != nil {
		return err
	}
	return uc.classrooms.Delete(ctx, req.OrgID, req.ClassroomID)
}
