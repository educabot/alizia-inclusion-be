package inclusion

import (
	"context"
	"errors"
	"strings"
	"time"

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

func (r *adaptationRepo) List(ctx context.Context, orgID uuid.UUID, filter providers.AdaptationFilter) ([]entities.Adaptation, error) {
	var adaptations []entities.Adaptation
	q := r.db.WithContext(ctx).
		Preload("Student").
		Preload("Teacher").
		Preload("Device").
		Preload("Devices").
		Where("organization_id = ?", orgID)
	if filter.TeacherID != nil {
		q = q.Where("teacher_id = ?", *filter.TeacherID)
	}
	if filter.StudentID != nil {
		q = q.Where("student_id = ?", *filter.StudentID)
	}
	if filter.DeviceID != nil {
		// material de valija usado: device principal o cualquiera del m2m.
		q = q.Where("device_id = ? OR id IN (SELECT adaptation_id FROM adaptation_devices WHERE device_id = ?)",
			*filter.DeviceID, *filter.DeviceID)
	}
	if q2 := strings.TrimSpace(filter.Query); q2 != "" {
		like := "%" + q2 + "%"
		q = q.Where("subject ILIKE ? OR title ILIKE ? OR coalesce(activity_description, '') ILIKE ?", like, like, like)
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
		Preload("Devices").
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

func (r *adaptationRepo) SetDevices(ctx context.Context, adaptationID int64, deviceIDs []int64) error {
	adaptation := &entities.Adaptation{ID: adaptationID}
	devices := make([]entities.Device, len(deviceIDs))
	for i, id := range deviceIDs {
		devices[i] = entities.Device{ID: id}
	}
	return r.db.WithContext(ctx).Model(adaptation).Association("Devices").Replace(devices)
}

func (r *adaptationRepo) CountSince(ctx context.Context, orgID uuid.UUID, since time.Time) (int, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&entities.Adaptation{}).
		Where("organization_id = ? AND created_at >= ?", orgID, since).
		Count(&count).Error
	return int(count), err
}

func (r *adaptationRepo) TopDevices(ctx context.Context, orgID uuid.UUID, limit int) ([]providers.DeviceUsageStat, error) {
	var results []providers.DeviceUsageStat
	err := r.db.WithContext(ctx).
		Table("adaptation_devices ad").
		Select("ad.device_id, d.name as device_name, COUNT(*) as count").
		Joins("JOIN devices d ON d.id = ad.device_id").
		Joins("JOIN adaptations a ON a.id = ad.adaptation_id").
		Where("a.organization_id = ?", orgID).
		Group("ad.device_id, d.name").
		Order("count DESC").
		Limit(limit).
		Find(&results).Error
	return results, err
}
