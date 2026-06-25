package entities

import (
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

// TeacherProfile holds teacher context (1:1 with users). All fields beyond the
// FKs are optional: dual age granularity (age_range and birthdate) plus
// pedagogical data that enriches the LLM prompt but is not required.
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
