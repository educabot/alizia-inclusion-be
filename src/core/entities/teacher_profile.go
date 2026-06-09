package entities

import (
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

// TeacherProfile es el contexto del docente (1:1 con users). Todos los campos
// más allá de las FKs son opcionales: doble granularidad de edad (age_range y
// birthdate) y datos pedagógicos que enriquecen el prompt sin ser obligatorios.
type TeacherProfile struct {
	ID              int64          `json:"id" gorm:"primaryKey"`
	UserID          int64          `json:"user_id"`
	OrganizationID  uuid.UUID      `json:"organization_id"`
	Birthdate       *time.Time     `json:"birthdate,omitempty"`
	AgeRange        *string        `json:"age_range,omitempty"`
	YearsExperience *int           `json:"years_experience,omitempty"`
	Specialization  *string        `json:"specialization,omitempty"`
	Subjects        pq.StringArray `json:"subjects,omitempty" gorm:"type:text[]"`
	TonePreference  *string        `json:"tone_preference,omitempty"`
	Bio             *string        `json:"bio,omitempty"`
	TimeTrackedEntity
}

func (TeacherProfile) TableName() string {
	return "teacher_profiles"
}
