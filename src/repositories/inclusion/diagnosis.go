package inclusion

import (
	"context"

	"gorm.io/gorm"

	"github.com/educabot/alizia-inclusion-be/src/core/entities"
	"github.com/educabot/alizia-inclusion-be/src/core/providers"
)

type diagnosisRepo struct {
	db *gorm.DB
}

func NewDiagnosisRepo(db *gorm.DB) providers.DiagnosisProvider {
	return &diagnosisRepo{db: db}
}

// ListByStudentProfile devuelve los diagnósticos de un perfil, con la etiqueta
// del catálogo precargada. El filtrado multi-tenant ya está garantizado por la
// cadena student → student_profile que resolvió el caller.
func (r *diagnosisRepo) ListByStudentProfile(ctx context.Context, studentProfileID int64) ([]entities.StudentDiagnosis, error) {
	var diagnoses []entities.StudentDiagnosis
	err := r.db.WithContext(ctx).
		Preload("Diagnosis").
		Where("student_profile_id = ?", studentProfileID).
		Order("id ASC").
		Find(&diagnoses).Error
	if err != nil {
		return nil, err
	}
	return diagnoses, nil
}
