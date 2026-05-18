package management

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/educabot/alizia-inclusion-be/src/core/entities"
	"github.com/educabot/alizia-inclusion-be/src/core/providers"
)

type classroomRepo struct {
	db *gorm.DB
}

func NewClassroomRepo(db *gorm.DB) providers.ClassroomProvider {
	return &classroomRepo{db: db}
}

func (r *classroomRepo) List(ctx context.Context, orgID uuid.UUID) ([]entities.Classroom, error) {
	var classrooms []entities.Classroom
	err := r.db.WithContext(ctx).
		Preload("Students").
		Where("organization_id = ?", orgID).
		Order("name ASC").
		Find(&classrooms).Error
	if err != nil {
		return nil, err
	}
	return classrooms, nil
}

func (r *classroomRepo) Get(ctx context.Context, orgID uuid.UUID, id int64) (*entities.Classroom, error) {
	var classroom entities.Classroom
	err := r.db.WithContext(ctx).
		Preload("Students").
		Where("organization_id = ? AND id = ?", orgID, id).
		First(&classroom).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, providers.ErrNotFound
		}
		return nil, err
	}
	return &classroom, nil
}

func (r *classroomRepo) Create(ctx context.Context, classroom *entities.Classroom) error {
	return r.db.WithContext(ctx).Create(classroom).Error
}

func (r *classroomRepo) Update(ctx context.Context, classroom *entities.Classroom) error {
	return r.db.WithContext(ctx).Save(classroom).Error
}

func (r *classroomRepo) Delete(ctx context.Context, orgID uuid.UUID, id int64) error {
	result := r.db.WithContext(ctx).
		Where("organization_id = ? AND id = ?", orgID, id).
		Delete(&entities.Classroom{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return providers.ErrNotFound
	}
	return nil
}
