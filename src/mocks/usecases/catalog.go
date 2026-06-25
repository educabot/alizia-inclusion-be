package mocks

import (
	"context"

	"github.com/stretchr/testify/mock"

	"github.com/educabot/alizia-inclusion-be/src/core/entities"
	cataloguc "github.com/educabot/alizia-inclusion-be/src/core/usecases/catalog"
)

type MockListRamps struct {
	mock.Mock
}

func (m *MockListRamps) Execute(ctx context.Context, req cataloguc.ListRampsRequest) ([]entities.Ramp, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]entities.Ramp), args.Error(1)
}

type MockGetRamp struct {
	mock.Mock
}

func (m *MockGetRamp) Execute(ctx context.Context, req cataloguc.GetRampRequest) (*entities.Ramp, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.Ramp), args.Error(1)
}

type MockListDevices struct {
	mock.Mock
}

func (m *MockListDevices) Execute(ctx context.Context, req cataloguc.ListDevicesRequest) ([]entities.Device, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]entities.Device), args.Error(1)
}

type MockGetDevice struct {
	mock.Mock
}

func (m *MockGetDevice) Execute(ctx context.Context, req cataloguc.GetDeviceRequest) (*entities.Device, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.Device), args.Error(1)
}
