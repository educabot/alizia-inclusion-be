package inclusion

import (
	"context"
	"errors"
	"strconv"
	"strings"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"gorm.io/gorm"

	"github.com/educabot/alizia-inclusion-be/src/core/entities"
	"github.com/educabot/alizia-inclusion-be/src/core/providers"
)

type pedagogicalContentRepo struct {
	db *gorm.DB
}

func NewPedagogicalContentRepo(db *gorm.DB) providers.PedagogicalContentProvider {
	return &pedagogicalContentRepo{db: db}
}

// searchDocument builds the weighted tsvector: title and keywords carry weight A,
// tags weight B, and chunk body weight C — keyword matches rank above body matches.
const searchDocument = `
	setweight(to_tsvector('spanish', coalesce(pc.title, '')), 'A') ||
	setweight(to_tsvector('spanish', array_to_string(pc.keywords, ' ')), 'A') ||
	setweight(to_tsvector('spanish', array_to_string(c.tags, ' ')), 'B') ||
	setweight(to_tsvector('spanish', coalesce(c.chunk_text, '')), 'C')`

// orTSQuery converts the query to an OR-semantics tsquery (any keyword matches,
// ranked by how many hit) — the expected RAG behavior. plainto_tsquery joins
// terms with AND; we replace ' & ' with ' | ' to get OR semantics.
const orTSQuery = `replace(plainto_tsquery('spanish', ?)::text, ' & ', ' | ')::tsquery`

// previewMaxChars caps the chunk preview returned by the RAG search.
const previewMaxChars = 280

// searchRow mirrors the SELECT columns; keywords is scanned as pq.StringArray
// and mapped to []string in the provider result.
type searchRow struct {
	ContentID int64
	ChunkID   int64
	Title     string
	Type      string
	Keywords  pq.StringArray `gorm:"type:text[]"`
	Preview   string
	Score     float64
}

func (r *pedagogicalContentRepo) SearchChunks(ctx context.Context, orgID uuid.UUID, query string, limit int) ([]providers.ContentSearchResult, error) {
	query = strings.TrimSpace(query)
	if query == "" || limit <= 0 {
		return []providers.ContentSearchResult{}, nil
	}

	sql := `
		SELECT pc.id AS content_id,
		       c.id AS chunk_id,
		       coalesce(pc.title, '') AS title,
		       coalesce(pc.type, '') AS type,
		       pc.keywords AS keywords,
		       left(coalesce(c.chunk_text, ''), ` + strconv.Itoa(previewMaxChars) + `) AS preview,
		       ts_rank((` + searchDocument + `), ` + orTSQuery + `) AS score
		FROM pedagogical_content_chunks c
		JOIN pedagogical_content pc ON pc.id = c.content_id
		WHERE (pc.organization_id IS NULL OR pc.organization_id = ?)
		  AND pc.status = 'published'
		  AND (` + searchDocument + `) @@ ` + orTSQuery + `
		ORDER BY score DESC, c.id ASC
		LIMIT ?`

	var rows []searchRow
	if err := r.db.WithContext(ctx).Raw(sql, query, orgID, query, limit).Scan(&rows).Error; err != nil {
		return nil, err
	}

	results := make([]providers.ContentSearchResult, len(rows))
	for i := range rows {
		results[i] = providers.ContentSearchResult{
			ContentID: rows[i].ContentID,
			ChunkID:   rows[i].ChunkID,
			Title:     rows[i].Title,
			Type:      rows[i].Type,
			Keywords:  rows[i].Keywords,
			Preview:   rows[i].Preview,
			Score:     rows[i].Score,
		}
	}
	return results, nil
}

func (r *pedagogicalContentRepo) GetContent(ctx context.Context, orgID uuid.UUID, contentID int64) (*entities.PedagogicalContent, error) {
	var content entities.PedagogicalContent
	err := r.db.WithContext(ctx).
		Preload("Chunks").
		Where("(organization_id IS NULL OR organization_id = ?) AND id = ?", orgID, contentID).
		First(&content).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, providers.ErrNotFound
		}
		return nil, err
	}
	return &content, nil
}
