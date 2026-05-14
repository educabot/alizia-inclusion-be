package inclusion

import (
	"context"

	"github.com/google/uuid"

	"github.com/educabot/alizia-inclusion-be/src/core/entities"
	"github.com/educabot/alizia-inclusion-be/src/core/providers"
)

type UpsertStudentProfileRequest struct {
	OrgID           uuid.UUID
	StudentID       int64
	IsTransitory    bool
	Difficulties    []string
	FreeDescription *string
}

func (r UpsertStudentProfileRequest) Validate() error {
	if r.OrgID == uuid.Nil {
		return errOrgIDRequired
	}
	if r.StudentID <= 0 {
		return errStudentIDRequired
	}
	return nil
}

type UpsertStudentProfile interface {
	Execute(ctx context.Context, req UpsertStudentProfileRequest) (*entities.StudentProfile, error)
}

type upsertStudentProfileImpl struct {
	students providers.StudentProvider
	profiles providers.StudentProfileProvider
}

func NewUpsertStudentProfile(students providers.StudentProvider, profiles providers.StudentProfileProvider) UpsertStudentProfile {
	return &upsertStudentProfileImpl{students: students, profiles: profiles}
}

func (uc *upsertStudentProfileImpl) Execute(ctx context.Context, req UpsertStudentProfileRequest) (*entities.StudentProfile, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	if _, err := uc.students.GetStudent(ctx, req.OrgID, req.StudentID); err != nil {
		return nil, err
	}

	profile := &entities.StudentProfile{
		StudentID:       req.StudentID,
		IsTransitory:    req.IsTransitory,
		Difficulties:    req.Difficulties,
		FreeDescription: req.FreeDescription,
	}

	if err := uc.profiles.Upsert(ctx, profile); err != nil {
		return nil, err
	}

	return uc.profiles.GetByStudentID(ctx, req.StudentID)
}
