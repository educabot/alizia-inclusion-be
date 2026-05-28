package mocks

import (
	"context"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"

	"github.com/educabot/alizia-inclusion-be/src/core/entities"
)

type MockDeviceProvider struct {
	mock.Mock
}

func (m *MockDeviceProvider) ListDevices(ctx context.Context, orgID uuid.UUID, rampID *int64) ([]entities.Device, error) {
	args := m.Called(ctx, orgID, rampID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]entities.Device), args.Error(1)
}

func (m *MockDeviceProvider) GetDevice(ctx context.Context, orgID uuid.UUID, id int64) (*entities.Device, error) {
	args := m.Called(ctx, orgID, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.Device), args.Error(1)
}
