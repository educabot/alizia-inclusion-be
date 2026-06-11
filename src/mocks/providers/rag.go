package mocks

import (
	"context"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"

	"github.com/educabot/alizia-inclusion-be/src/core/entities"
	"github.com/educabot/alizia-inclusion-be/src/core/providers"
)

type MockPedagogicalContentProvider struct {
	mock.Mock
}

func (m *MockPedagogicalContentProvider) SearchChunks(ctx context.Context, orgID uuid.UUID, query string, limit int) ([]providers.ContentSearchResult, error) {
	args := m.Called(ctx, orgID, query, limit)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]providers.ContentSearchResult), args.Error(1)
}

func (m *MockPedagogicalContentProvider) GetContent(ctx context.Context, orgID uuid.UUID, contentID int64) (*entities.PedagogicalContent, error) {
	args := m.Called(ctx, orgID, contentID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.PedagogicalContent), args.Error(1)
}
