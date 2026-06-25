package entities

import (
	"time"

	"github.com/google/uuid"
)

// Diagnosis is a catalog entry for clinical/educational diagnoses. It is a
// secondary, optional layer relative to observable situations.
// organization_id NULL means the record is global (Educabot-managed);
// a non-null value means it belongs to a specific organization.
type Diagnosis struct {
	ID             int64      `json:"id" gorm:"primaryKey"`
	OrganizationID *uuid.UUID `json:"organization_id,omitempty"`
	Name           string     `json:"name"`
	Category       *string    `json:"category,omitempty"`
}

func (Diagnosis) TableName() string {
	return "diagnoses_catalog"
}

// StudentDiagnosis links a student profile to a diagnosis catalog entry.
// This structured layer is only populated when the school provides a formal diagnosis.
type StudentDiagnosis struct {
	ID               int64      `json:"id" gorm:"primaryKey"`
	StudentProfileID int64      `json:"student_profile_id"`
	DiagnosisID      int64      `json:"diagnosis_id"`
	Severity         *string    `json:"severity,omitempty"`
	Notes            *string    `json:"notes,omitempty"`
	Diagnosis        *Diagnosis `json:"diagnosis,omitempty" gorm:"foreignKey:DiagnosisID"`
	CreatedAt        time.Time  `json:"created_at"`
}

func (StudentDiagnosis) TableName() string {
	return "student_diagnoses"
}
