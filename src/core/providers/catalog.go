package providers

import (
	"context"

	"github.com/google/uuid"

	"github.com/educabot/alizia-inclusion-be/src/core/entities"
)

type RampProvider interface {
	ListRamps(ctx context.Context, orgID uuid.UUID) ([]entities.Ramp, error)
	GetRamp(ctx context.Context, orgID uuid.UUID, id int64) (*entities.Ramp, error)
}

type DeviceProvider interface {
	ListDevices(ctx context.Context, orgID uuid.UUID, rampID *int64) ([]entities.Device, error)
	GetDevice(ctx context.Context, orgID uuid.UUID, id int64) (*entities.Device, error)
}
