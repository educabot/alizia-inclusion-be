package management

import (
	"context"

	"github.com/google/uuid"

	"github.com/educabot/alizia-inclusion-be/src/core/entities"
	"github.com/educabot/alizia-inclusion-be/src/core/providers"
)

type GetClassroomRequest struct {
	OrgID       uuid.UUID
	ClassroomID int64
}

func (r GetClassroomRequest) Validate() error {
	if r.OrgID == uuid.Nil {
		return errOrgIDRequired
	}
	if r.ClassroomID <= 0 {
		return errClassroomIDRequired
	}
	return nil
}

type GetClassroom interface {
	Execute(ctx context.Context, req GetClassroomRequest) (*entities.Classroom, error)
}

type getClassroomImpl struct {
	classrooms providers.ClassroomProvider
}

func NewGetClassroom(classrooms providers.ClassroomProvider) GetClassroom {
	return &getClassroomImpl{classrooms: classrooms}
}

func (uc *getClassroomImpl) Execute(ctx context.Context, req GetClassroomRequest) (*entities.Classroom, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	return uc.classrooms.Get(ctx, req.OrgID, req.ClassroomID)
}
