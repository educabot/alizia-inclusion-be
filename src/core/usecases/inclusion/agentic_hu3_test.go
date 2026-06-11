package inclusion

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/educabot/alizia-inclusion-be/src/core/entities"
	"github.com/educabot/alizia-inclusion-be/src/core/providers"
	mockproviders "github.com/educabot/alizia-inclusion-be/src/mocks/providers"
)

func TestInclusionDispatcher_SearchContentReturnsRankedResults(t *testing.T) {
	// Arrange
	ctx := context.Background()
	orgID := uuid.New()
	content := new(mockproviders.MockPedagogicalContentProvider)
	content.On("SearchChunks", ctx, orgID, "TEA autismo", defaultContentSearchLimit).
		Return([]providers.ContentSearchResult{
			{ContentID: 2, ChunkID: 2, Title: "Estrategias TEA", Preview: "anticipación…", Score: 0.46},
		}, nil)
	d := inclusionDispatcher{content: content}

	// Act
	result, err := d.Dispatch(ctx, orgID, providers.ToolCall{
		Name:      "search_content",
		Arguments: `{"query": "TEA autismo"}`,
	})

	// Assert
	require.NoError(t, err)
	var got struct {
		Results []providers.ContentSearchResult `json:"results"`
	}
	require.NoError(t, json.Unmarshal([]byte(result), &got))
	require.Len(t, got.Results, 1)
	assert.Equal(t, "Estrategias TEA", got.Results[0].Title)
	content.AssertExpectations(t)
}

func TestInclusionDispatcher_SearchContentNoMatchReturnsEmpty(t *testing.T) {
	// Arrange
	ctx := context.Background()
	orgID := uuid.New()
	content := new(mockproviders.MockPedagogicalContentProvider)
	content.On("SearchChunks", ctx, orgID, "quimica organica", defaultContentSearchLimit).
		Return([]providers.ContentSearchResult{}, nil)
	d := inclusionDispatcher{content: content}

	// Act
	result, err := d.Dispatch(ctx, orgID, providers.ToolCall{
		Name:      "search_content",
		Arguments: `{"query": "quimica organica"}`,
	})

	// Assert
	require.NoError(t, err)
	var got struct {
		Results []providers.ContentSearchResult `json:"results"`
	}
	require.NoError(t, json.Unmarshal([]byte(result), &got))
	assert.Empty(t, got.Results)
}

func TestInclusionDispatcher_GetContentReturnsDocument(t *testing.T) {
	// Arrange
	ctx := context.Background()
	orgID := uuid.New()
	title := "Estrategias TEA"
	content := new(mockproviders.MockPedagogicalContentProvider)
	content.On("GetContent", ctx, orgID, int64(2)).
		Return(&entities.PedagogicalContent{ID: 2, Title: &title, Status: "published"}, nil)
	d := inclusionDispatcher{content: content}

	// Act
	result, err := d.Dispatch(ctx, orgID, providers.ToolCall{
		Name:      "get_content",
		Arguments: `{"content_id": 2}`,
	})

	// Assert
	require.NoError(t, err)
	var got entities.PedagogicalContent
	require.NoError(t, json.Unmarshal([]byte(result), &got))
	assert.Equal(t, int64(2), got.ID)
	content.AssertExpectations(t)
}

func TestInclusionDispatcher_ContentToolsUnavailableWhenProviderNil(t *testing.T) {
	// Arrange
	d := inclusionDispatcher{}

	// Act
	_, errSearch := d.Dispatch(context.Background(), uuid.New(), providers.ToolCall{
		Name:      "search_content",
		Arguments: `{"query": "TEA"}`,
	})
	_, errGet := d.Dispatch(context.Background(), uuid.New(), providers.ToolCall{
		Name:      "get_content",
		Arguments: `{"content_id": 1}`,
	})

	// Assert
	require.Error(t, errSearch)
	require.Error(t, errGet)
	assert.Contains(t, errSearch.Error(), "no disponible")
}

func TestInclusionTools_ExposeSearchAndGetContent(t *testing.T) {
	// Arrange / Act
	names := make(map[string]bool)
	for _, tool := range inclusionTools() {
		names[tool.Name] = true
	}

	// Assert
	assert.True(t, names["search_content"])
	assert.True(t, names["get_content"])
}
