package entities

import (
	"time"

	"github.com/google/uuid"
)

// StudentNote es una nota de seguimiento del alumno (perfil). Acá viven el
// "diagnóstico" u observables del alta (internal=true, no visible en el front) y
// el seguimiento a lo largo del tiempo. Es el historial del alumno, NO de un
// recurso puntual: por eso vive en el perfil y no en la adaptación.
type StudentNote struct {
	ID             int64     `json:"id" gorm:"primaryKey"`
	StudentID      int64     `json:"student_id"`
	OrganizationID uuid.UUID `json:"organization_id"`
	// UserID es el docente dueño de la nota: las notas son privadas, cada docente
	// ve solo las suyas. Las filas legacy (user_id NULL) quedan invisibles.
	UserID    int64     `json:"user_id"`
	Content   string    `json:"content"`
	Type      string    `json:"type" gorm:"default:seguimiento"`
	Internal  bool      `json:"internal" gorm:"default:true"`
	CreatedAt time.Time `json:"created_at"`
}

func (StudentNote) TableName() string {
	return "student_notes"
}
