package mocks

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"

	"github.com/educabot/alizia-inclusion-be/src/core/providers"
)

type MockAIUsageProvider struct {
	mock.Mock
}

func (m *MockAIUsageProvider) Record(ctx context.Context, record providers.AIUsageRecord) error {
	args := m.Called(ctx, record)
	return args.Error(0)
}

func (m *MockAIUsageProvider) Summarize(ctx context.Context, orgID uuid.UUID, since time.Time) (*providers.AIUsageSummary, error) {
	args := m.Called(ctx, orgID, since)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*providers.AIUsageSummary), args.Error(1)
}
