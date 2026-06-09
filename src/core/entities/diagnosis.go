package entities

import (
	"time"

	"github.com/google/uuid"
)

// Diagnosis es una etiqueta del catálogo de diagnósticos (capa secundaria y
// opcional respecto de las situaciones observables). organization_id NULL =
// global (Educabot); con valor = propio de la organización.
type Diagnosis struct {
	ID             int64      `json:"id" gorm:"primaryKey"`
	OrganizationID *uuid.UUID `json:"organization_id,omitempty"`
	Name           string     `json:"name"`
	Category       *string    `json:"category,omitempty"`
}

func (Diagnosis) TableName() string {
	return "diagnoses_catalog"
}

// StudentDiagnosis liga un perfil de alumno con un diagnóstico del catálogo.
// Capa estructurada, solo se llena si la escuela brinda diagnóstico.
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
