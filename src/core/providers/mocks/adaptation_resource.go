package mocks

import (
	"context"

	"github.com/stretchr/testify/mock"

	"github.com/educabot/alizia-inclusion-be/src/core/entities"
)

type MockAdaptationResourceProvider struct {
	mock.Mock
}

func (m *MockAdaptationResourceProvider) ListByAdaptation(ctx context.Context, adaptationID int64) ([]entities.AdaptationResource, error) {
	args := m.Called(ctx, adaptationID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]entities.AdaptationResource), args.Error(1)
}
