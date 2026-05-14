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
