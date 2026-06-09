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
	// Campos enriquecidos (HU-2, todos opcionales): doble granularidad de edad
	// y nombre preferido para personalizar la respuesta sin exigir datos.
	Birthdate     *time.Time      `json:"birthdate,omitempty"`
	AgeRange      *string         `json:"age_range,omitempty"`
	GradeLevel    *string         `json:"grade_level,omitempty"`
	PreferredName *string         `json:"preferred_name,omitempty"`
	Profile       *StudentProfile `json:"profile,omitempty" gorm:"foreignKey:StudentID"`
	TimeTrackedEntity
}
