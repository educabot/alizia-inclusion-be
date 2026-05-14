package catalog

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/educabot/alizia-inclusion-be/src/core/entities"
	"github.com/educabot/alizia-inclusion-be/src/core/providers"
)

type rampRepo struct {
	db *gorm.DB
}

func NewRampRepo(db *gorm.DB) providers.RampProvider {
	return &rampRepo{db: db}
}

func (r *rampRepo) ListRamps(ctx context.Context, orgID uuid.UUID) ([]entities.Ramp, error) {
	var ramps []entities.Ramp
	err := r.db.WithContext(ctx).
		Preload("Devices", func(db *gorm.DB) *gorm.DB {
			return db.Order("sort_order ASC")
		}).
		Where("organization_id = ?", orgID).
		Order("sort_order ASC").
		Find(&ramps).Error
	if err != nil {
		return nil, err
	}
	return ramps, nil
}

func (r *rampRepo) GetRamp(ctx context.Context, orgID uuid.UUID, id int64) (*entities.Ramp, error) {
	var ramp entities.Ramp
	err := r.db.WithContext(ctx).
		Preload("Devices", func(db *gorm.DB) *gorm.DB {
			return db.Order("sort_order ASC")
		}).
		Where("organization_id = ? AND id = ?", orgID, id).
		First(&ramp).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, providers.ErrNotFound
		}
		return nil, err
	}
	return &ramp, nil
}
