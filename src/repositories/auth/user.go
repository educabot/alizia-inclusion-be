package auth

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/educabot/alizia-inclusion-be/src/core/entities"
	"github.com/educabot/alizia-inclusion-be/src/core/providers"
)

type userRepo struct {
	db *gorm.DB
}

func NewUserRepo(db *gorm.DB) providers.UserProvider {
	return &userRepo{db: db}
}

func (r *userRepo) GetByID(ctx context.Context, orgID uuid.UUID, id int64) (*entities.User, error) {
	var user entities.User
	err := r.db.WithContext(ctx).
		Where("organization_id = ? AND id = ?", orgID, id).
		First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, providers.ErrNotFound
		}
		return nil, err
	}
	return &user, nil
}

func (r *userRepo) ListByRole(ctx context.Context, orgID uuid.UUID, role string) ([]entities.User, error) {
	var users []entities.User
	err := r.db.WithContext(ctx).
		Where("organization_id = ? AND role = ?", orgID, role).
		Order("name ASC").
		Find(&users).Error
	if err != nil {
		return nil, err
	}
	return users, nil
}
