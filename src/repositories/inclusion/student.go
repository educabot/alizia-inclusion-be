package inclusion

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/educabot/alizia-inclusion-be/src/core/entities"
	"github.com/educabot/alizia-inclusion-be/src/core/providers"
)

type studentRepo struct {
	db *gorm.DB
}

func NewStudentRepo(db *gorm.DB) providers.StudentProvider {
	return &studentRepo{db: db}
}

func (r *studentRepo) GetStudent(ctx context.Context, orgID uuid.UUID, id int64) (*entities.Student, error) {
	var student entities.Student
	err := r.db.WithContext(ctx).
		Preload("Profile").
		Where("organization_id = ? AND id = ?", orgID, id).
		First(&student).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, providers.ErrNotFound
		}
		return nil, err
	}
	return &student, nil
}

func (r *studentRepo) ListByClassroom(ctx context.Context, orgID uuid.UUID, classroomID int64) ([]entities.Student, error) {
	var students []entities.Student
	err := r.db.WithContext(ctx).
		Preload("Profile").
		Where("organization_id = ? AND classroom_id = ?", orgID, classroomID).
		Order("name ASC").
		Find(&students).Error
	if err != nil {
		return nil, err
	}
	return students, nil
}

func (r *studentRepo) List(ctx context.Context, orgID uuid.UUID) ([]entities.Student, error) {
	var students []entities.Student
	err := r.db.WithContext(ctx).
		Preload("Profile").
		Where("organization_id = ?", orgID).
		Order("name ASC").
		Find(&students).Error
	if err != nil {
		return nil, err
	}
	return students, nil
}

func (r *studentRepo) Create(ctx context.Context, student *entities.Student) error {
	return r.db.WithContext(ctx).Create(student).Error
}

func (r *studentRepo) Update(ctx context.Context, student *entities.Student) error {
	return r.db.WithContext(ctx).Save(student).Error
}

func (r *studentRepo) Delete(ctx context.Context, orgID uuid.UUID, id int64) error {
	result := r.db.WithContext(ctx).
		Where("organization_id = ? AND id = ?", orgID, id).
		Delete(&entities.Student{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return providers.ErrNotFound
	}
	return nil
}
