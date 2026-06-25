package mocks

import (
	"context"

	"github.com/stretchr/testify/mock"

	"github.com/educabot/alizia-inclusion-be/src/core/providers"
)

type MockEmbedder struct {
	mock.Mock
}

func (m *MockEmbedder) EmbedQuery(ctx context.Context, text string) ([]float32, error) {
	args := m.Called(ctx, text)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]float32), args.Error(1)
}

type MockRAGSearchProvider struct {
	mock.Mock
}

func (m *MockRAGSearchProvider) HybridSearch(ctx context.Context, spec providers.HybridSearchSpec, embedding []float32) ([]providers.ChunkHit, error) {
	args := m.Called(ctx, spec, embedding)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]providers.ChunkHit), args.Error(1)
}
