package providers

import (
	"context"

	"github.com/google/uuid"

	"github.com/educabot/alizia-inclusion-be/src/core/entities"
)

// TeacherProfileProvider — contexto del docente (1:1 con users).
type TeacherProfileProvider interface {
	GetByUserID(ctx context.Context, orgID uuid.UUID, userID int64) (*entities.TeacherProfile, error)
}

// SituationCatalogProvider — catálogo de situaciones observables. Devuelve las
// globales (organization_id IS NULL) más las propias de la organización.
type SituationCatalogProvider interface {
	List(ctx context.Context, orgID uuid.UUID) ([]entities.Situation, error)
}

// DiagnosisProvider — diagnósticos estructurados ligados a un perfil de alumno.
type DiagnosisProvider interface {
	ListByStudentProfile(ctx context.Context, studentProfileID int64) ([]entities.StudentDiagnosis, error)
}

// PPIProvider — Proyecto Pedagógico Individual (1 por alumno).
type PPIProvider interface {
	GetByStudentID(ctx context.Context, orgID uuid.UUID, studentID int64) (*entities.PPI, error)
}
