package entities

import (
	"time"

	"github.com/google/uuid"
)

// IntegradoraAssignment liga una maestra integradora (user) con un alumno
// asignado. Habilita que el motor de contexto sepa qué alumnos puede trabajar
// la integradora. El RBAC queda fuera de scope (otro equipo).
type IntegradoraAssignment struct {
	ID             int64     `json:"id" gorm:"primaryKey"`
	OrganizationID uuid.UUID `json:"organization_id"`
	UserID         int64     `json:"user_id"`
	StudentID      int64     `json:"student_id"`
	CreatedAt      time.Time `json:"created_at"`
}

func (IntegradoraAssignment) TableName() string {
	return "integradora_assignments"
}
