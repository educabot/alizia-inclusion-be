package entities

import (
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

// PedagogicalContent is a document in the pedagogical corpus (book, paper,
// material, chapter). Notion-style hierarchy via self-referencing parent_id.
// Independent of device/ramp scope. organization_id NULL means global (Educabot).
// RAG retrieval is driven by keywords[] with a GIN index.
type PedagogicalContent struct {
	ID             int64                     `json:"id" gorm:"primaryKey"`
	ParentID       *int64                    `json:"parent_id,omitempty"`
	Type           *string                   `json:"type,omitempty"`
	Title          *string                   `json:"title,omitempty"`
	Status         string                    `json:"status"`
	Keywords       pq.StringArray            `json:"keywords,omitempty" gorm:"type:text[]"`
	OrganizationID *uuid.UUID                `json:"organization_id,omitempty"`
	Chunks         []PedagogicalContentChunk `json:"chunks,omitempty" gorm:"foreignKey:ContentID"`
	CreatedAt      time.Time                 `json:"created_at"`
	UpdatedAt      time.Time                 `json:"updated_at"`
}

func (PedagogicalContent) TableName() string {
	return "pedagogical_content"
}

// PedagogicalContentChunk is a searchable unit of a document. For the MVP,
// 1 chunk = the whole document; splitting logic changes when books arrive,
// no migration needed. The embedding column exists in the table but is inert
// for now (keyword/full-text first) and is therefore not modelled in Go yet.
type PedagogicalContentChunk struct {
	ID        int64          `json:"id" gorm:"primaryKey"`
	ContentID int64          `json:"content_id"`
	ChunkText *string        `json:"chunk_text,omitempty"`
	Tags      pq.StringArray `json:"tags,omitempty" gorm:"type:text[]"`
	CreatedAt time.Time      `json:"created_at"`
}

func (PedagogicalContentChunk) TableName() string {
	return "pedagogical_content_chunks"
}
