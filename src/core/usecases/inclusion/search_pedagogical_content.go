package inclusion

import (
	"context"

	"github.com/google/uuid"

	"github.com/educabot/alizia-inclusion-be/src/core/providers"
)

// defaultSearchContentLimit is the default top-N for RAG queries.
const defaultSearchContentLimit = 5

// maxSearchContentLimit caps the top-N a caller may request. Without a ceiling,
// an oversized LIMIT from the API would hit the DB unnecessarily for a RAG top-N.
const maxSearchContentLimit = 20

type SearchContentRequest struct {
	OrgID uuid.UUID
	Query string
	Limit int
}

func (r SearchContentRequest) Validate() error {
	if r.OrgID == uuid.Nil {
		return errOrgIDRequired
	}
	return nil
}

type SearchContentResponse struct {
	Query   string                          `json:"query"`
	Results []providers.ContentSearchResult `json:"results"`
}

// SearchPedagogicalContent exposes the RAG keyword/full-text search without
// invoking the LLM. Used for corpus validation and as the backing implementation
// of the search_content tool in the agentic loop.
type SearchPedagogicalContent interface {
	Execute(ctx context.Context, req SearchContentRequest) (*SearchContentResponse, error)
}

type searchPedagogicalContentImpl struct {
	content providers.PedagogicalContentProvider
}

func NewSearchPedagogicalContent(content providers.PedagogicalContentProvider) SearchPedagogicalContent {
	return &searchPedagogicalContentImpl{content: content}
}

func (uc *searchPedagogicalContentImpl) Execute(ctx context.Context, req SearchContentRequest) (*SearchContentResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	limit := req.Limit
	if limit <= 0 {
		limit = defaultSearchContentLimit
	}
	if limit > maxSearchContentLimit {
		limit = maxSearchContentLimit
	}

	results, err := uc.content.SearchChunks(ctx, req.OrgID, req.Query, limit)
	if err != nil {
		return nil, err
	}
	// Ensure non-nil slice so JSON serializes as [] rather than null on empty results.
	if results == nil {
		results = []providers.ContentSearchResult{}
	}
	return &SearchContentResponse{Query: req.Query, Results: results}, nil
}
