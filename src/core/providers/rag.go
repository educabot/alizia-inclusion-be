package providers

import (
	"context"

	"github.com/google/uuid"

	"github.com/educabot/alizia-inclusion-be/src/core/entities"
)

// ContentSearchResult is a RAG-retrieved chunk with its preview and a
// relevance score (higher = more relevant). Carries parent-document metadata
// so the LLM can decide whether to fetch the full content via get_content.
type ContentSearchResult struct {
	ContentID int64    `json:"content_id"`
	ChunkID   int64    `json:"chunk_id"`
	Title     string   `json:"title"`
	Type      string   `json:"type,omitempty"`
	Keywords  []string `json:"keywords,omitempty"`
	Preview   string   `json:"preview"`
	Score     float64  `json:"score"`
}

// PedagogicalContentProvider is the RAG interface for pedagogical content (Layer E).
// Keyword/full-text search in MVP; embedding-based search is deferred.
type PedagogicalContentProvider interface {
	// SearchChunks returns the most relevant chunks for a keyword/full-text query,
	// ordered by descending score. Returns an empty slice (not an error) when
	// nothing matches — the LLM must fall back to base guidelines, not hallucinate.
	SearchChunks(ctx context.Context, orgID uuid.UUID, query string, limit int) ([]ContentSearchResult, error)
	// GetContent retrieves a document together with all its chunks.
	GetContent(ctx context.Context, orgID uuid.UUID, contentID int64) (*entities.PedagogicalContent, error)
}
