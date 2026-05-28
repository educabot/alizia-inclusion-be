package mocks

import (
	"context"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"

	"github.com/educabot/alizia-inclusion-be/src/core/entities"
)

type MockRampProvider struct {
	mock.Mock
}

func (m *MockRampProvider) ListRamps(ctx context.Context, orgID uuid.UUID) ([]entities.Ramp, error) {
	args := m.Called(ctx, orgID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]entities.Ramp), args.Error(1)
}

func (m *MockRampProvider) GetRamp(ctx context.Context, orgID uuid.UUID, id int64) (*entities.Ramp, error) {
	args := m.Called(ctx, orgID, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.Ramp), args.Error(1)
}
