package inclusion_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/educabot/alizia-inclusion-be/src/core/providers"
	mockproviders "github.com/educabot/alizia-inclusion-be/src/mocks/providers"
	"github.com/educabot/alizia-inclusion-be/src/core/usecases/inclusion"
	"github.com/educabot/alizia-inclusion-be/src/testutil"
)

func TestSearchContent_RejectsNilOrgID(t *testing.T) {
	// Arrange
	content := new(mockproviders.MockPedagogicalContentProvider)
	uc := inclusion.NewSearchPedagogicalContent(content)

	// Act
	_, err := uc.Execute(context.Background(), inclusion.SearchContentRequest{Query: "TEA"})

	// Assert
	require.Error(t, err)
	content.AssertNotCalled(t, "SearchChunks")
}

func TestSearchContent_AppliesDefaultLimitAndReturnsResults(t *testing.T) {
	// Arrange
	content := new(mockproviders.MockPedagogicalContentProvider)
	content.On("SearchChunks", mock.Anything, testutil.TestOrgID, "TEA autismo", 5).
		Return([]providers.ContentSearchResult{
			{ContentID: 2, ChunkID: 2, Title: "Estrategias TEA", Preview: "…", Score: 0.46},
		}, nil)
	uc := inclusion.NewSearchPedagogicalContent(content)

	// Act
	got, err := uc.Execute(context.Background(), inclusion.SearchContentRequest{
		OrgID: testutil.TestOrgID,
		Query: "TEA autismo",
	})

	// Assert
	require.NoError(t, err)
	assert.Equal(t, "TEA autismo", got.Query)
	require.Len(t, got.Results, 1)
	assert.Equal(t, int64(2), got.Results[0].ContentID)
	content.AssertExpectations(t)
}

func TestSearchContent_ClampsExcessiveLimitToMax(t *testing.T) {
	// Arrange: caller requests an excessively large top-N; the usecase clamps it to 20
	// before hitting the DB (an unbounded LIMIT would hammer the RAG index pointlessly).
	content := new(mockproviders.MockPedagogicalContentProvider)
	content.On("SearchChunks", mock.Anything, testutil.TestOrgID, "TEA", 20).
		Return([]providers.ContentSearchResult{}, nil)
	uc := inclusion.NewSearchPedagogicalContent(content)

	// Act
	_, err := uc.Execute(context.Background(), inclusion.SearchContentRequest{
		OrgID: testutil.TestOrgID,
		Query: "TEA",
		Limit: 100000,
	})

	// Assert
	require.NoError(t, err)
	content.AssertExpectations(t)
}

func TestSearchContent_NoMatchReturnsEmptySliceNotNil(t *testing.T) {
	// Arrange
	content := new(mockproviders.MockPedagogicalContentProvider)
	content.On("SearchChunks", mock.Anything, testutil.TestOrgID, "quimica organica", 3).
		Return([]providers.ContentSearchResult(nil), nil)
	uc := inclusion.NewSearchPedagogicalContent(content)

	// Act
	got, err := uc.Execute(context.Background(), inclusion.SearchContentRequest{
		OrgID: testutil.TestOrgID,
		Query: "quimica organica",
		Limit: 3,
	})

	// Assert
	require.NoError(t, err)
	require.NotNil(t, got.Results)
	assert.Empty(t, got.Results)
}
