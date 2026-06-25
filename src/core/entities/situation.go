package entities

import "github.com/google/uuid"

// Situation es una situación observable de aula (las ~15 del MVP: "no inicia la
// tarea", "se distrae constantemente", etc.). Es la entrada pedagógica primaria
// (se parte de lo observable, no del diagnóstico). organization_id NULL = global
// (catálogo de Educabot); con valor = definición propia de la organización.
type Situation struct {
	ID             int64      `json:"id" gorm:"primaryKey"`
	OrganizationID *uuid.UUID `json:"organization_id,omitempty"`
	Code           string     `json:"code"`
	Name           string     `json:"name"`
	Description    *string    `json:"description,omitempty"`
	Phase          *string    `json:"phase,omitempty"`
	SortOrder      int        `json:"sort_order"`
}

func (Situation) TableName() string {
	return "situations_catalog"
}
