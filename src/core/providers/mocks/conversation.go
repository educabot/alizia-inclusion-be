package mocks

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"

	"github.com/educabot/alizia-inclusion-be/src/core/entities"
	"github.com/educabot/alizia-inclusion-be/src/core/providers"
)

type MockConversationProvider struct {
	mock.Mock
}

func (m *MockConversationProvider) ListByUser(ctx context.Context, orgID uuid.UUID, userID int64, mode string) ([]entities.Conversation, error) {
	args := m.Called(ctx, orgID, userID, mode)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]entities.Conversation), args.Error(1)
}

func (m *MockConversationProvider) AppendTurn(ctx context.Context, params providers.AppendTurnParams) (providers.AppendTurnResult, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(providers.AppendTurnResult), args.Error(1)
}

func (m *MockConversationProvider) Delete(ctx context.Context, orgID uuid.UUID, id int64) error {
	args := m.Called(ctx, orgID, id)
	return args.Error(0)
}

func (m *MockConversationProvider) Rename(ctx context.Context, orgID uuid.UUID, id int64, title string) error {
	args := m.Called(ctx, orgID, id, title)
	return args.Error(0)
}

func (m *MockConversationProvider) ListPendingSummary(ctx context.Context, idleBefore time.Time, limit int) ([]entities.Conversation, error) {
	args := m.Called(ctx, idleBefore, limit)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]entities.Conversation), args.Error(1)
}
