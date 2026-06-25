package entities

import (
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

// PPI (Proyecto Pedagógico Individual, 1 por alumno): objetivos, adaptaciones
// curriculares y seguimiento. Lo crea el docente y lo valida el Director. Es
// contexto de primera línea cuando existe; todos los campos opcionales.
type PPI struct {
	ID                    int64          `json:"id" gorm:"primaryKey"`
	OrganizationID        uuid.UUID      `json:"organization_id"`
	StudentID             int64          `json:"student_id"`
	Objectives            pq.StringArray `json:"objectives,omitempty" gorm:"type:text[]"`
	CurricularAdaptations *string        `json:"curricular_adaptations,omitempty"`
	FollowUp              *string        `json:"follow_up,omitempty"`
	Status                string         `json:"status"`
	CreatedBy             *int64         `json:"created_by,omitempty"`
	ValidatedBy           *int64         `json:"validated_by,omitempty"`
	CreatedAt             time.Time      `json:"created_at"`
	UpdatedAt             time.Time      `json:"updated_at"`
}

func (PPI) TableName() string {
	return "ppi"
}
