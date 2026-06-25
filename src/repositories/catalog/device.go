package catalog

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/educabot/alizia-inclusion-be/src/core/entities"
	"github.com/educabot/alizia-inclusion-be/src/core/providers"
)

type deviceRepo struct {
	db *gorm.DB
}

func NewDeviceRepo(db *gorm.DB) providers.DeviceProvider {
	return &deviceRepo{db: db}
}

func (r *deviceRepo) ListDevices(ctx context.Context, orgID uuid.UUID, rampID *int64) ([]entities.Device, error) {
	var devices []entities.Device
	q := r.db.WithContext(ctx).
		Preload("Ramp").
		Preload("Resources").
		Where("organization_id = ?", orgID).
		Where("is_active = ?", true)

	if rampID != nil {
		q = q.Where("ramp_id = ?", *rampID)
	}

	err := q.Order("sort_order ASC").Find(&devices).Error
	if err != nil {
		return nil, err
	}
	return devices, nil
}

func (r *deviceRepo) GetDevice(ctx context.Context, orgID uuid.UUID, id int64) (*entities.Device, error) {
	var device entities.Device
	err := r.db.WithContext(ctx).
		Preload("Ramp").
		Preload("Resources").
		Where("organization_id = ? AND id = ?", orgID, id).
		First(&device).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, providers.ErrNotFound
		}
		return nil, err
	}
	return &device, nil
}
