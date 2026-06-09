package inclusion

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/educabot/alizia-inclusion-be/src/core/entities"
	"github.com/educabot/alizia-inclusion-be/src/core/providers"
)

type teacherProfileRepo struct {
	db *gorm.DB
}

func NewTeacherProfileRepo(db *gorm.DB) providers.TeacherProfileProvider {
	return &teacherProfileRepo{db: db}
}

func (r *teacherProfileRepo) GetByUserID(ctx context.Context, orgID uuid.UUID, userID int64) (*entities.TeacherProfile, error) {
	var profile entities.TeacherProfile
	err := r.db.WithContext(ctx).
		Where("organization_id = ? AND user_id = ?", orgID, userID).
		First(&profile).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, providers.ErrNotFound
		}
		return nil, err
	}
	return &profile, nil
}
