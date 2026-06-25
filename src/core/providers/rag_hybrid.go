package providers

import "context"

// Embedder genera el vector de embedding de un texto para la búsqueda semántica.
// Se implementa contra el endpoint de embeddings de Azure (recurso propio, ver config).
type Embedder interface {
	EmbedQuery(ctx context.Context, text string) ([]float32, error)
}

// HybridSearchSpec son los parámetros de una búsqueda híbrida sobre el corpus rag_*.
// SemanticQuestion alimenta el embedding; Terms refuerza FTS de content/summary,
// términos exactos y conceptos. ResourceID (opcional) acota a un documento.
type HybridSearchSpec struct {
	ResourceID       *int64
	SemanticQuestion string
	Terms            []string
	Limit            int
	Offset           int
}

// ChunkHit es un fragmento de rag_chunks con su score híbrido y la metadata del
// recurso padre (rag_resources). Es el resultado que consume la tool/endpoint.
type ChunkHit struct {
	ChunkID    int64    `json:"chunk_id"`
	ResourceID int64    `json:"resource_id"`
	Title      string   `json:"title"`
	ChunkIndex int32    `json:"chunk_index"`
	PageStart  int32    `json:"page_start"`
	PageEnd    int32    `json:"page_end"`
	Score      float64  `json:"score"`
	Sources    string   `json:"sources"`
	Summary    string   `json:"summary,omitempty"`
	Concepts   []string `json:"concepts,omitempty"`
	Content    string   `json:"content"`
}

// RAGSearchProvider — búsqueda híbrida (vector + FTS de content/summary + términos
// exactos + conceptos) sobre el corpus global rag_resources/rag_chunks. El embedding
// de la pregunta lo provee el caller (Embedder); el provider solo corre el ranking SQL.
type RAGSearchProvider interface {
	HybridSearch(ctx context.Context, spec HybridSearchSpec, embedding []float32) ([]ChunkHit, error)
}
