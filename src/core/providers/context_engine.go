package providers

import (
	"context"

	"github.com/google/uuid"

	"github.com/educabot/alizia-inclusion-be/src/core/entities"
)

// TeacherProfileProvider retrieves teacher profile context (1:1 with users).
type TeacherProfileProvider interface {
	GetByUserID(ctx context.Context, orgID uuid.UUID, userID int64) (*entities.TeacherProfile, error)
}

// SituationCatalogProvider lists observable situations: global ones
// (organization_id IS NULL) plus those belonging to the organization.
type SituationCatalogProvider interface {
	List(ctx context.Context, orgID uuid.UUID) ([]entities.Situation, error)
}

// DiagnosisProvider retrieves structured diagnoses linked to a student profile.
type DiagnosisProvider interface {
	ListByStudentProfile(ctx context.Context, studentProfileID int64) ([]entities.StudentDiagnosis, error)
}

// PPIProvider retrieves the Individual Pedagogical Project (PPI) for a student (1 per student).
type PPIProvider interface {
	GetByStudentID(ctx context.Context, orgID uuid.UUID, studentID int64) (*entities.PPI, error)
}

// IntegradoraAssignmentProvider manages the integration-teacher ↔ student assignment.
type IntegradoraAssignmentProvider interface {
	ListStudentIDsByUser(ctx context.Context, orgID uuid.UUID, userID int64) ([]int64, error)
	IsAssigned(ctx context.Context, orgID uuid.UUID, userID, studentID int64) (bool, error)
}
