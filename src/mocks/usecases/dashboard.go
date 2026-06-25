package mocks

import (
	"context"

	"github.com/stretchr/testify/mock"

	dashuc "github.com/educabot/alizia-inclusion-be/src/core/usecases/dashboard"
)

type MockGetMetrics struct {
	mock.Mock
}

func (m *MockGetMetrics) Execute(ctx context.Context, req dashuc.GetMetricsRequest) (*dashuc.GetMetricsResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dashuc.GetMetricsResponse), args.Error(1)
}

type MockGetAIUsage struct {
	mock.Mock
}

func (m *MockGetAIUsage) Execute(ctx context.Context, req dashuc.GetAIUsageRequest) (*dashuc.GetAIUsageResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dashuc.GetAIUsageResponse), args.Error(1)
}
