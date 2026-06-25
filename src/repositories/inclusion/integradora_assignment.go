package inclusion

import (
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/educabot/alizia-inclusion-be/src/core/entities"
	"github.com/educabot/alizia-inclusion-be/src/core/providers"
)

type integradoraAssignmentRepo struct {
	db *gorm.DB
}

func NewIntegradoraAssignmentRepo(db *gorm.DB) providers.IntegradoraAssignmentProvider {
	return &integradoraAssignmentRepo{db: db}
}

func (r *integradoraAssignmentRepo) ListStudentIDsByUser(ctx context.Context, orgID uuid.UUID, userID int64) ([]int64, error) {
	var ids []int64
	err := r.db.WithContext(ctx).
		Model(&entities.IntegradoraAssignment{}).
		Where("organization_id = ? AND user_id = ?", orgID, userID).
		Order("student_id ASC").
		Pluck("student_id", &ids).Error
	if err != nil {
		return nil, err
	}
	return ids, nil
}

func (r *integradoraAssignmentRepo) IsAssigned(ctx context.Context, orgID uuid.UUID, userID, studentID int64) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&entities.IntegradoraAssignment{}).
		Where("organization_id = ? AND user_id = ? AND student_id = ?", orgID, userID, studentID).
		Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
