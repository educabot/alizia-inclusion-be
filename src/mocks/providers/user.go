package mocks

import (
	"context"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"

	"github.com/educabot/alizia-inclusion-be/src/core/entities"
)

type MockUserProvider struct {
	mock.Mock
}

func (m *MockUserProvider) GetByID(ctx context.Context, orgID uuid.UUID, id int64) (*entities.User, error) {
	args := m.Called(ctx, orgID, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.User), args.Error(1)
}

func (m *MockUserProvider) ListByRole(ctx context.Context, orgID uuid.UUID, role string) ([]entities.User, error) {
	args := m.Called(ctx, orgID, role)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]entities.User), args.Error(1)
}
