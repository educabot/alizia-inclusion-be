package inclusion

import (
	"context"

	"gorm.io/gorm"

	"github.com/educabot/alizia-inclusion-be/src/core/entities"
	"github.com/educabot/alizia-inclusion-be/src/core/providers"
)

type adaptationResourceRepo struct {
	db *gorm.DB
}

func NewAdaptationResourceRepo(db *gorm.DB) providers.AdaptationResourceProvider {
	return &adaptationResourceRepo{db: db}
}

func (r *adaptationResourceRepo) ListByAdaptation(ctx context.Context, adaptationID int64) ([]entities.AdaptationResource, error) {
	var resources []entities.AdaptationResource
	err := r.db.WithContext(ctx).
		Where("adaptation_id = ?", adaptationID).
		Order("created_at DESC").
		Find(&resources).Error
	return resources, err
}
