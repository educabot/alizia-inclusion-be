package inclusion

import (
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/educabot/alizia-inclusion-be/src/core/entities"
	"github.com/educabot/alizia-inclusion-be/src/core/providers"
)

type situationRepo struct {
	db *gorm.DB
}

func NewSituationRepo(db *gorm.DB) providers.SituationCatalogProvider {
	return &situationRepo{db: db}
}

// List devuelve las situaciones globales (organization_id IS NULL) más las
// propias de la organización, ordenadas por sort_order.
func (r *situationRepo) List(ctx context.Context, orgID uuid.UUID) ([]entities.Situation, error) {
	var situations []entities.Situation
	err := r.db.WithContext(ctx).
		Where("organization_id IS NULL OR organization_id = ?", orgID).
		Order("sort_order ASC, code ASC").
		Find(&situations).Error
	if err != nil {
		return nil, err
	}
	return situations, nil
}
