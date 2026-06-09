package providers

import (
	"context"

	"github.com/google/uuid"

	"github.com/educabot/alizia-inclusion-be/src/core/entities"
)

// ContentSearchResult es un chunk encontrado por el RAG, con su preview y un
// score de relevancia (mayor = más pertinente). Lleva metadata del documento
// padre para que la LLM decida si profundizar con get_content.
type ContentSearchResult struct {
	ContentID int64    `json:"content_id"`
	ChunkID   int64    `json:"chunk_id"`
	Title     string   `json:"title"`
	Type      string   `json:"type,omitempty"`
	Keywords  []string `json:"keywords,omitempty"`
	Preview   string   `json:"preview"`
	Score     float64  `json:"score"`
}

// PedagogicalContentProvider — RAG de contenido pedagógico (Capa E). Búsqueda
// keyword/full-text en el MVP; embeddings quedan inertes (Futuro).
type PedagogicalContentProvider interface {
	// SearchChunks busca por keywords/full-text y devuelve los chunks más
	// relevantes primero. Sin coincidencias devuelve slice vacío (no error):
	// la LLM debe caer a los lineamientos base sin inventar.
	SearchChunks(ctx context.Context, orgID uuid.UUID, query string, limit int) ([]ContentSearchResult, error)
	// GetContent trae un documento con sus chunks (chunk/documento completo).
	GetContent(ctx context.Context, orgID uuid.UUID, contentID int64) (*entities.PedagogicalContent, error)
}
