package auth

import (
	"context"

	"github.com/google/uuid"

	"github.com/educabot/alizia-inclusion-be/src/core/entities"
	"github.com/educabot/alizia-inclusion-be/src/core/providers"
)

type GetMeRequest struct {
	OrgID  uuid.UUID
	UserID int64
}

func (r GetMeRequest) Validate() error {
	if r.OrgID == uuid.Nil {
		return errOrgIDRequired
	}
	if r.UserID <= 0 {
		return errUserIDRequired
	}
	return nil
}

type GetMe interface {
	Execute(ctx context.Context, req GetMeRequest) (*entities.User, error)
}

type getMeImpl struct {
	users providers.UserProvider
}

func NewGetMe(users providers.UserProvider) GetMe {
	return &getMeImpl{users: users}
}

func (uc *getMeImpl) Execute(ctx context.Context, req GetMeRequest) (*entities.User, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	return uc.users.GetByID(ctx, req.OrgID, req.UserID)
}
