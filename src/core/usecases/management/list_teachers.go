package management

import (
	"context"

	"github.com/google/uuid"

	"github.com/educabot/alizia-inclusion-be/src/core/entities"
	"github.com/educabot/alizia-inclusion-be/src/core/providers"
)

type ListTeachersRequest struct {
	OrgID uuid.UUID
}

func (r ListTeachersRequest) Validate() error {
	if r.OrgID == uuid.Nil {
		return errOrgIDRequired
	}
	return nil
}

type ListTeachers interface {
	Execute(ctx context.Context, req ListTeachersRequest) ([]entities.User, error)
}

type listTeachersImpl struct {
	users providers.UserProvider
}

func NewListTeachers(users providers.UserProvider) ListTeachers {
	return &listTeachersImpl{users: users}
}

func (uc *listTeachersImpl) Execute(ctx context.Context, req ListTeachersRequest) ([]entities.User, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	return uc.users.ListByRole(ctx, req.OrgID, "teacher")
}
