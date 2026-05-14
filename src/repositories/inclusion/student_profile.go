package inclusion

import (
	"context"
	"errors"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/educabot/alizia-inclusion-be/src/core/entities"
	"github.com/educabot/alizia-inclusion-be/src/core/providers"
)

type studentProfileRepo struct {
	db *gorm.DB
}

func NewStudentProfileRepo(db *gorm.DB) providers.StudentProfileProvider {
	return &studentProfileRepo{db: db}
}

func (r *studentProfileRepo) GetByStudentID(ctx context.Context, studentID int64) (*entities.StudentProfile, error) {
	var profile entities.StudentProfile
	err := r.db.WithContext(ctx).
		Where("student_id = ?", studentID).
		First(&profile).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, providers.ErrProfileNotFound
		}
		return nil, err
	}
	return &profile, nil
}

func (r *studentProfileRepo) Upsert(ctx context.Context, profile *entities.StudentProfile) error {
	return r.db.WithContext(ctx).
		Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "student_id"}},
			DoUpdates: clause.AssignmentColumns([]string{"is_transitory", "difficulties", "free_description", "updated_at"}),
		}).
		Create(profile).Error
}
