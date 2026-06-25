package inclusion

import (
	"context"
	"fmt"
	"strings"

	"github.com/educabot/alizia-inclusion-be/src/core/providers"
)

const defaultHybridSearchLimit = 5

// HybridSearchRequest son los parámetros de la búsqueda híbrida sobre el corpus rag_*.
// La pregunta completa alimenta el embedding; Terms refuerza FTS/términos/conceptos.
type HybridSearchRequest struct {
	SemanticQuestion string
	Terms            []string
	ResourceID       *int64
	Limit            int
	Offset           int
}

func (r HybridSearchRequest) Validate() error {
	if strings.TrimSpace(r.SemanticQuestion) == "" {
		return errQuestionRequired
	}
	return nil
}

type HybridSearchResponse struct {
	Question string                `json:"question"`
	Results  []providers.ChunkHit `json:"results"`
}

// HybridSearchContent embebe la pregunta y corre la búsqueda híbrida RAG. Es la base
// del endpoint (validación por Postman) y de la tool agéntica search_content_hibrido.
type HybridSearchContent interface {
	Execute(ctx context.Context, req HybridSearchRequest) (*HybridSearchResponse, error)
}

type hybridSearchContentImpl struct {
	embedder providers.Embedder
	rag      providers.RAGSearchProvider
}

func NewHybridSearchContent(embedder providers.Embedder, rag providers.RAGSearchProvider) HybridSearchContent {
	return &hybridSearchContentImpl{embedder: embedder, rag: rag}
}

func (uc *hybridSearchContentImpl) Execute(ctx context.Context, req HybridSearchRequest) (*HybridSearchResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	limit := req.Limit
	if limit <= 0 {
		limit = defaultHybridSearchLimit
	}

	embedding, err := uc.embedder.EmbedQuery(ctx, req.SemanticQuestion)
	if err != nil {
		return nil, fmt.Errorf("embed query: %w", err)
	}

	hits, err := uc.rag.HybridSearch(ctx, providers.HybridSearchSpec{
		ResourceID:       req.ResourceID,
		SemanticQuestion: req.SemanticQuestion,
		Terms:            req.Terms,
		Limit:            limit,
		Offset:           req.Offset,
	}, embedding)
	if err != nil {
		return nil, err
	}
	// Slice no-nil para que el JSON sea [] y no null cuando no hay match.
	if hits == nil {
		hits = []providers.ChunkHit{}
	}
	return &HybridSearchResponse{Question: req.SemanticQuestion, Results: hits}, nil
}
