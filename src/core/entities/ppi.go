package entities

import (
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

// PPI (Individual Pedagogical Project) holds one record per student with learning
// objectives, curricular adaptations, and follow-up notes. Created by the teacher,
// validated by the Director. Treated as first-priority context when present; all
// content fields are optional.
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
