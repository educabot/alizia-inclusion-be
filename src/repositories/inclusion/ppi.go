package inclusion

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/educabot/alizia-inclusion-be/src/core/entities"
	"github.com/educabot/alizia-inclusion-be/src/core/providers"
)

type ppiRepo struct {
	db *gorm.DB
}

func NewPPIRepo(db *gorm.DB) providers.PPIProvider {
	return &ppiRepo{db: db}
}

func (r *ppiRepo) GetByStudentID(ctx context.Context, orgID uuid.UUID, studentID int64) (*entities.PPI, error) {
	var ppi entities.PPI
	err := r.db.WithContext(ctx).
		Where("organization_id = ? AND student_id = ?", orgID, studentID).
		First(&ppi).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, providers.ErrNotFound
		}
		return nil, err
	}
	return &ppi, nil
}
