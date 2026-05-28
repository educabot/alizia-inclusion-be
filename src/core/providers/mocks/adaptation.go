package mocks

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"

	"github.com/educabot/alizia-inclusion-be/src/core/entities"
	"github.com/educabot/alizia-inclusion-be/src/core/providers"
)

type MockAdaptationProvider struct {
	mock.Mock
}

func (m *MockAdaptationProvider) List(ctx context.Context, orgID uuid.UUID, studentID *int64) ([]entities.Adaptation, error) {
	args := m.Called(ctx, orgID, studentID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]entities.Adaptation), args.Error(1)
}

func (m *MockAdaptationProvider) Get(ctx context.Context, orgID uuid.UUID, id int64) (*entities.Adaptation, error) {
	args := m.Called(ctx, orgID, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.Adaptation), args.Error(1)
}

func (m *MockAdaptationProvider) Create(ctx context.Context, adaptation *entities.Adaptation) error {
	args := m.Called(ctx, adaptation)
	return args.Error(0)
}

func (m *MockAdaptationProvider) Update(ctx context.Context, adaptation *entities.Adaptation) error {
	args := m.Called(ctx, adaptation)
	return args.Error(0)
}

func (m *MockAdaptationProvider) Delete(ctx context.Context, orgID uuid.UUID, id int64) error {
	args := m.Called(ctx, orgID, id)
	return args.Error(0)
}

func (m *MockAdaptationProvider) SetDevices(ctx context.Context, adaptationID int64, deviceIDs []int64) error {
	args := m.Called(ctx, adaptationID, deviceIDs)
	return args.Error(0)
}

func (m *MockAdaptationProvider) CountSince(ctx context.Context, orgID uuid.UUID, since time.Time) (int, error) {
	args := m.Called(ctx, orgID, since)
	return args.Int(0), args.Error(1)
}

func (m *MockAdaptationProvider) TopDevices(ctx context.Context, orgID uuid.UUID, limit int) ([]providers.DeviceUsageStat, error) {
	args := m.Called(ctx, orgID, limit)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]providers.DeviceUsageStat), args.Error(1)
}
