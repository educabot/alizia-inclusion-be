package entities

import (
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

// PedagogicalContent es un documento del corpus pedagógico (libro / paper /
// material / capítulo). Jerárquico tipo Notion (parent_id self-ref). Es
// independiente de la valija: no tiene device_id ni ramp_id. organization_id
// NULL = global (Educabot). El RAG busca por keywords[] (GIN).
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

// PedagogicalContentChunk es un pedacito buscable de un documento. Para el MVP
// 1 chunk = el documento entero; cuando lleguen libros solo cambia el paso de
// partir en chunks (sin migración). La columna embedding existe en la tabla
// pero queda inerte en el MVP (keyword/full-text first), por eso no se modela
// en Go todavía.
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
