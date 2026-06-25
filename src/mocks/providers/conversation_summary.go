package mocks

import (
	"context"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"

	"github.com/educabot/alizia-inclusion-be/src/core/entities"
)

type MockConversationSummaryProvider struct {
	mock.Mock
}

func (m *MockConversationSummaryProvider) RecentByStudent(ctx context.Context, orgID uuid.UUID, studentID int64, limit int) ([]entities.ConversationSummary, error) {
	args := m.Called(ctx, orgID, studentID, limit)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]entities.ConversationSummary), args.Error(1)
}

func (m *MockConversationSummaryProvider) RecentByDevice(ctx context.Context, orgID uuid.UUID, deviceID int64, limit int) ([]entities.ConversationSummary, error) {
	args := m.Called(ctx, orgID, deviceID, limit)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]entities.ConversationSummary), args.Error(1)
}

func (m *MockConversationSummaryProvider) RecentByTopic(ctx context.Context, orgID uuid.UUID, keyword string, limit int) ([]entities.ConversationSummary, error) {
	args := m.Called(ctx, orgID, keyword, limit)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]entities.ConversationSummary), args.Error(1)
}

func (m *MockConversationSummaryProvider) Upsert(ctx context.Context, summary *entities.ConversationSummary, studentIDs, deviceIDs []int64) error {
	args := m.Called(ctx, summary, studentIDs, deviceIDs)
	return args.Error(0)
}
