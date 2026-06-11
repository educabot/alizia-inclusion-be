package inclusion

import (
	"context"

	"github.com/google/uuid"

	"github.com/educabot/alizia-inclusion-be/src/core/providers"
)

// defaultSearchContentLimit acota el top-N por defecto del RAG.
const defaultSearchContentLimit = 5

// maxSearchContentLimit es el techo del top-N que un caller puede pedir. Acota el
// trabajo del RAG aunque el request pida un Limit desmedido (la query viene de la
// API; sin techo, un LIMIT enorme golpearía la DB sin sentido para un top-N).
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

// SearchPedagogicalContent expone el buscador del RAG (keyword/full-text) sin
// pasar por la LLM. Sirve para validar el corpus por Postman y como base de la
// tool search_content del loop agéntico.
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
	// Garantizamos slice no-nil para que el JSON sea [] y no null cuando no hay match.
	if results == nil {
		results = []providers.ContentSearchResult{}
	}
	return &SearchContentResponse{Query: req.Query, Results: results}, nil
}
