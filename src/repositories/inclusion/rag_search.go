package inclusion

import (
	"context"
	"fmt"
	"log/slog"
	"strings"

	"github.com/lib/pq"
	"gorm.io/gorm"

	"github.com/educabot/alizia-inclusion-be/src/core/providers"
)

type ragSearchRepo struct {
	db *gorm.DB
}

func NewRAGSearchRepo(db *gorm.DB) providers.RAGSearchProvider {
	return &ragSearchRepo{db: db}
}

// ragChunkRow espeja el SELECT final; concepts entra como text[] vía pq.StringArray.
type ragChunkRow struct {
	ChunkID    int64
	ResourceID int64
	Title      string
	ChunkIndex int32
	PageStart  int32
	PageEnd    int32
	Score      float64
	Sources    string
	Summary    string
	Concepts   pq.StringArray `gorm:"type:text[]"`
	Content    string
}

// HybridSearch combina vector (60%) + FTS de content (18%) y summary (12%) + términos
// exactos (6%) + conceptos (8%) en un único ranking sobre rag_chunks/rag_resources.
func (r *ragSearchRepo) HybridSearch(ctx context.Context, spec providers.HybridSearchSpec, embedding []float32) ([]providers.ChunkHit, error) {
	limit := spec.Limit
	if limit <= 0 {
		limit = 5
	}
	offset := spec.Offset
	if offset < 0 {
		offset = 0
	}
	// Margen de candidatos por sub-consulta antes del merge/paginado final.
	candidateLimit := offset + limit + 80

	var resourceID any
	if spec.ResourceID != nil {
		resourceID = *spec.ResourceID
	}

	var rows []ragChunkRow
	// lib/pq no soporta @name. GORM convierte ? → $N para PostgreSQL, por eso usamos ?
	// (si usamos $N directamente GORM cuenta 0 bindvars y manda 0 args → error).
	// @content_query y @summary_query aparecen 2 veces cada uno en el CTE → 9 args total.
	querySQL := hybridSearchSQL
	querySQL = strings.ReplaceAll(querySQL, "__QVEC__", pgVector(embedding))
	querySQL = strings.ReplaceAll(querySQL, "@content_query", "?")
	querySQL = strings.ReplaceAll(querySQL, "@summary_query", "?")
	querySQL = strings.ReplaceAll(querySQL, "@terms", "?")
	querySQL = strings.ReplaceAll(querySQL, "@resource_id", "?")
	querySQL = strings.ReplaceAll(querySQL, "@candidate_limit", "?")
	querySQL = strings.ReplaceAll(querySQL, "@result_limit", "?")
	querySQL = strings.ReplaceAll(querySQL, "@result_offset", "?")

	contentQuery := webSearchQuery(spec.Terms)
	summaryQuery := strings.Join(spec.Terms, " ")
	err := r.db.WithContext(ctx).Raw(querySQL,
		contentQuery,               // ? content_query (valor)
		contentQuery,               // ? content_query (websearch_to_tsquery)
		summaryQuery,               // ? summary_query (valor)
		summaryQuery,               // ? summary_query (websearch_to_tsquery)
		pq.StringArray(spec.Terms), // ? terms
		resourceID,                 // ? resource_id
		candidateLimit,             // ? candidate_limit
		limit,                      // ? result_limit
		offset,                     // ? result_offset
	).Scan(&rows).Error
	if err != nil {
		return nil, err
	}

	var topScore float64
	if len(rows) > 0 {
		topScore = rows[0].Score
	}
	slog.InfoContext(ctx, "rag.search",
		"terms", spec.Terms,
		"resource_id", resourceID,
		"limit", limit,
		"offset", offset,
		"candidate_limit", candidateLimit,
		"hits", len(rows),
		"top_score", topScore,
	)

	hits := make([]providers.ChunkHit, len(rows))
	for i := range rows {
		hits[i] = providers.ChunkHit{
			ChunkID:    rows[i].ChunkID,
			ResourceID: rows[i].ResourceID,
			Title:      rows[i].Title,
			ChunkIndex: rows[i].ChunkIndex,
			PageStart:  rows[i].PageStart,
			PageEnd:    rows[i].PageEnd,
			Score:      rows[i].Score,
			Sources:    rows[i].Sources,
			Summary:    rows[i].Summary,
			Concepts:   rows[i].Concepts,
			Content:    rows[i].Content,
		}
	}
	return hits, nil
}

// pgVector serializa el embedding al literal de pgvector ("[v1,v2,...]").
func pgVector(v []float32) string {
	parts := make([]string, len(v))
	for i, value := range v {
		parts[i] = fmt.Sprintf("%.8f", value)
	}
	return "[" + strings.Join(parts, ",") + "]"
}

// webSearchQuery arma la query para websearch_to_tsquery: términos con espacios o
// guiones van entre comillas (frase), unidos por OR.
func webSearchQuery(terms []string) string {
	var parts []string
	for _, term := range terms {
		term = strings.TrimSpace(term)
		if term == "" {
			continue
		}
		term = strings.ReplaceAll(term, `"`, `\"`)
		if strings.ContainsAny(term, " -") {
			term = `"` + term + `"`
		}
		parts = append(parts, term)
	}
	return strings.Join(parts, " OR ")
}

// hybridSearchSQL: ranking híbrido sobre rag_chunks/rag_resources.
// Parámetros: $1=content_query $2=summary_query $3=terms $4=resource_id
//             $5=candidate_limit $6=result_limit $7=result_offset
// El vector se inyecta inline vía strings.ReplaceAll(__QVEC__) antes de ejecutar.
const hybridSearchSQL = `
WITH params AS (
  SELECT
    '__QVEC__'::vector AS qvec,
    @content_query::text AS content_query,
    websearch_to_tsquery('pg_catalog.simple'::regconfig, @content_query::text) AS content_tsq,
    @summary_query::text AS summary_query,
    websearch_to_tsquery('pg_catalog.simple'::regconfig, @summary_query::text) AS summary_tsq,
    @terms::text[] AS terms,
    @resource_id::bigint AS resource_id,
    @candidate_limit::int AS candidate_limit,
    @result_limit::int AS result_limit,
    @result_offset::int AS result_offset
),

vector_hits AS (
  SELECT
    c.id,
    'vector' AS source,
    0.60 * (1 - (c.embedding <=> p.qvec)) AS weighted_score
  FROM rag_chunks c, params p
  WHERE p.resource_id IS NULL OR c.resource_id = p.resource_id
  ORDER BY c.embedding <=> p.qvec, c.id
  LIMIT (SELECT candidate_limit FROM params)
),

content_fts_hits AS (
  SELECT
    c.id,
    'content_fts' AS source,
    0.18 * LEAST(1.0, ts_rank_cd(
      to_tsvector('pg_catalog.simple'::regconfig, COALESCE(c.content, '')),
      p.content_tsq
    ) * 8) AS weighted_score
  FROM rag_chunks c, params p
  WHERE p.content_query <> ''
    AND (p.resource_id IS NULL OR c.resource_id = p.resource_id)
    AND to_tsvector('pg_catalog.simple'::regconfig, COALESCE(c.content, '')) @@ p.content_tsq
  ORDER BY weighted_score DESC, c.id
  LIMIT (SELECT candidate_limit FROM params)
),

summary_fts_hits AS (
  SELECT
    c.id,
    'summary_fts' AS source,
    0.12 * LEAST(1.0, ts_rank_cd(
      to_tsvector('pg_catalog.simple'::regconfig, COALESCE(c.summary, '')),
      p.summary_tsq
    ) * 8) AS weighted_score
  FROM rag_chunks c, params p
  WHERE p.summary_query <> ''
    AND (p.resource_id IS NULL OR c.resource_id = p.resource_id)
    AND to_tsvector('pg_catalog.simple'::regconfig, COALESCE(c.summary, '')) @@ p.summary_tsq
  ORDER BY weighted_score DESC, c.id
  LIMIT (SELECT candidate_limit FROM params)
),

exact_term_hits AS (
  SELECT
    c.id,
    'exact_terms' AS source,
    0.06 * LEAST(
      1.0,
      COUNT(DISTINCT term.term)::float / GREATEST(cardinality(p.terms), 1)
    ) AS weighted_score
  FROM rag_chunks c
  CROSS JOIN params p
  CROSS JOIN LATERAL unnest(p.terms) AS term(term)
  WHERE (p.resource_id IS NULL OR c.resource_id = p.resource_id)
    AND (
      lower(c.content) LIKE '%' || lower(term.term) || '%'
      OR lower(COALESCE(c.summary, '')) LIKE '%' || lower(term.term) || '%'
    )
  GROUP BY c.id, p.terms
  ORDER BY weighted_score DESC, c.id
  LIMIT (SELECT candidate_limit FROM params)
),

concept_hits AS (
  SELECT
    c.id,
    'concepts' AS source,
    0.08 * LEAST(
      1.0,
      COUNT(DISTINCT qconcept.term)::float / GREATEST(cardinality(p.terms), 1)
    ) AS weighted_score
  FROM rag_chunks c
  CROSS JOIN params p
  CROSS JOIN LATERAL unnest(p.terms) AS qconcept(term)
  WHERE (p.resource_id IS NULL OR c.resource_id = p.resource_id)
    AND EXISTS (
      SELECT 1
      FROM unnest(c.concepts) AS cconcept(term)
      WHERE lower(cconcept.term) LIKE '%' || lower(qconcept.term) || '%'
         OR lower(qconcept.term) LIKE '%' || lower(cconcept.term) || '%'
    )
  GROUP BY c.id, p.terms
  ORDER BY weighted_score DESC, c.id
  LIMIT (SELECT candidate_limit FROM params)
),

all_hits AS (
  SELECT * FROM vector_hits
  UNION ALL SELECT * FROM content_fts_hits
  UNION ALL SELECT * FROM summary_fts_hits
  UNION ALL SELECT * FROM exact_term_hits
  UNION ALL SELECT * FROM concept_hits
),

ranked AS (
  SELECT
    id,
    SUM(weighted_score) AS score,
    string_agg(source, ',' ORDER BY source) AS sources
  FROM all_hits
  GROUP BY id
)

SELECT
  c.id AS chunk_id,
  c.resource_id AS resource_id,
  r.title AS title,
  c.chunk_index AS chunk_index,
  c.page_start AS page_start,
  c.page_end AS page_end,
  ranked.score AS score,
  ranked.sources AS sources,
  COALESCE(c.summary, '') AS summary,
  c.concepts AS concepts,
  c.content AS content
FROM ranked
JOIN rag_chunks c ON c.id = ranked.id
JOIN rag_resources r ON r.id = c.resource_id
ORDER BY ranked.score DESC, c.id
LIMIT (SELECT result_limit FROM params)
OFFSET (SELECT result_offset FROM params);
`
