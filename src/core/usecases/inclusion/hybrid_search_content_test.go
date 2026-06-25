package inclusion_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/educabot/alizia-inclusion-be/src/core/providers"
	mockproviders "github.com/educabot/alizia-inclusion-be/src/core/providers/mocks"
	"github.com/educabot/alizia-inclusion-be/src/core/usecases/inclusion"
)

func TestHybridSearchContent_RequiresSemanticQuestion(t *testing.T) {
	// Arrange
	uc := inclusion.NewHybridSearchContent(new(mockproviders.MockEmbedder), new(mockproviders.MockRAGSearchProvider))

	// Act
	_, err := uc.Execute(context.Background(), inclusion.HybridSearchRequest{SemanticQuestion: "   "})

	// Assert
	assert.ErrorIs(t, err, providers.ErrValidation)
}

func TestHybridSearchContent_EmbedsQuestionAndReturnsRankedHits(t *testing.T) {
	// Arrange
	ctx := context.Background()
	embedder := new(mockproviders.MockEmbedder)
	rag := new(mockproviders.MockRAGSearchProvider)
	emb := []float32{0.1, 0.2, 0.3}
	embedder.On("EmbedQuery", ctx, "buenas practicas dislexia?").Return(emb, nil)
	hits := []providers.ChunkHit{{ChunkID: 1, ResourceID: 3, Title: "Guia", Score: 0.91, Content: "texto"}}
	rag.On("HybridSearch", ctx, mock.AnythingOfType("providers.HybridSearchSpec"), emb).Return(hits, nil)

	// Act
	res, err := inclusion.NewHybridSearchContent(embedder, rag).Execute(ctx, inclusion.HybridSearchRequest{
		SemanticQuestion: "buenas practicas dislexia?",
		Terms:            []string{"dislexia"},
	})

	// Assert
	require.NoError(t, err)
	require.Len(t, res.Results, 1)
	assert.Equal(t, int64(3), res.Results[0].ResourceID)
	embedder.AssertExpectations(t)
	rag.AssertExpectations(t)
}

func TestHybridSearchContent_PropagatesEmbedErrorWithoutSearching(t *testing.T) {
	// Arrange
	ctx := context.Background()
	embedder := new(mockproviders.MockEmbedder)
	rag := new(mockproviders.MockRAGSearchProvider)
	embedder.On("EmbedQuery", ctx, "q").Return(nil, errors.New("azure down"))

	// Act
	_, err := inclusion.NewHybridSearchContent(embedder, rag).Execute(ctx, inclusion.HybridSearchRequest{SemanticQuestion: "q"})

	// Assert
	require.Error(t, err)
	rag.AssertNotCalled(t, "HybridSearch", mock.Anything, mock.Anything, mock.Anything)
}
