package auth

import (
	"context"
	"fmt"

	"golang.org/x/crypto/bcrypt"

	"github.com/educabot/alizia-inclusion-be/src/core/entities"
	"github.com/educabot/alizia-inclusion-be/src/core/providers"
)

type LoginRequest struct {
	Email    string
	Password string
}

func (r LoginRequest) Validate() error {
	if r.Email == "" {
		return errEmailRequired
	}
	if r.Password == "" {
		return errPasswordRequired
	}
	return nil
}

type Login interface {
	Execute(ctx context.Context, req LoginRequest) (*entities.User, error)
}

type loginImpl struct {
	users providers.UserProvider
}

func NewLogin(users providers.UserProvider) Login {
	return &loginImpl{users: users}
}

func (uc *loginImpl) Execute(ctx context.Context, req LoginRequest) (*entities.User, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	user, err := uc.users.GetByEmail(ctx, req.Email)
	if err != nil {
		return nil, fmt.Errorf("%w", providers.ErrInvalidCredentials)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		return nil, fmt.Errorf("%w", providers.ErrInvalidCredentials)
	}

	return user, nil
}
