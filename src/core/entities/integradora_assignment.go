package entities

import (
	"time"

	"github.com/google/uuid"
)

// IntegradoraAssignment links an integration teacher (user) to an assigned
// student, allowing the context engine to know which students the teacher
// may work with. RBAC enforcement is out of scope for this entity.
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
