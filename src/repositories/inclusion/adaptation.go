package inclusion

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/educabot/alizia-inclusion-be/src/core/entities"
	"github.com/educabot/alizia-inclusion-be/src/core/providers"
)

type adaptationRepo struct {
	db *gorm.DB
}

func NewAdaptationRepo(db *gorm.DB) providers.AdaptationProvider {
	return &adaptationRepo{db: db}
}

func (r *adaptationRepo) List(ctx context.Context, orgID uuid.UUID, studentID *int64) ([]entities.Adaptation, error) {
	var adaptations []entities.Adaptation
	q := r.db.WithContext(ctx).
		Preload("Student").
		Preload("Teacher").
		Preload("Device").
		Where("organization_id = ?", orgID)
	if studentID != nil {
		q = q.Where("student_id = ?", *studentID)
	}
	err := q.Order("created_at DESC").Find(&adaptations).Error
	if err != nil {
		return nil, err
	}
	return adaptations, nil
}

func (r *adaptationRepo) Get(ctx context.Context, orgID uuid.UUID, id int64) (*entities.Adaptation, error) {
	var adaptation entities.Adaptation
	err := r.db.WithContext(ctx).
		Preload("Student").
		Preload("Teacher").
		Preload("Device").
		Where("organization_id = ? AND id = ?", orgID, id).
		First(&adaptation).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, providers.ErrAdaptationNotFound
		}
		return nil, err
	}
	return &adaptation, nil
}

func (r *adaptationRepo) Create(ctx context.Context, adaptation *entities.Adaptation) error {
	return r.db.WithContext(ctx).Create(adaptation).Error
}

func (r *adaptationRepo) Update(ctx context.Context, adaptation *entities.Adaptation) error {
	return r.db.WithContext(ctx).Save(adaptation).Error
}

func (r *adaptationRepo) Delete(ctx context.Context, orgID uuid.UUID, id int64) error {
	result := r.db.WithContext(ctx).
		Where("organization_id = ? AND id = ?", orgID, id).
		Delete(&entities.Adaptation{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return providers.ErrAdaptationNotFound
	}
	return nil
}
