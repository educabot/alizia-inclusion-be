package entities

import (
	"time"

	"github.com/google/uuid"
)

type Student struct {
	ID             int64     `json:"id" gorm:"primaryKey"`
	OrganizationID uuid.UUID `json:"organization_id"`
	ClassroomID    int64     `json:"classroom_id"`
	Name           string    `json:"name"`
	// Enriched fields (HU-2, all optional): dual age granularity and preferred name
	// to personalise AI responses without requiring complete data.
	Birthdate     *time.Time      `json:"birthdate,omitempty"`
	AgeRange      *string         `json:"age_range,omitempty"`
	GradeLevel    *string         `json:"grade_level,omitempty"`
	PreferredName *string         `json:"preferred_name,omitempty"`
	Profile       *StudentProfile `json:"profile,omitempty" gorm:"foreignKey:StudentID"`
	TimeTrackedEntity
}
