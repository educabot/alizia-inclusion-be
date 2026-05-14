package catalog

import (
	"context"

	"github.com/google/uuid"

	"github.com/educabot/alizia-inclusion-be/src/core/entities"
	"github.com/educabot/alizia-inclusion-be/src/core/providers"
)

type ListDevicesRequest struct {
	OrgID  uuid.UUID
	RampID *int64
}

func (r ListDevicesRequest) Validate() error {
	if r.OrgID == uuid.Nil {
		return errOrgIDRequired
	}
	return nil
}

type ListDevices interface {
	Execute(ctx context.Context, req ListDevicesRequest) ([]entities.Device, error)
}

type listDevicesImpl struct {
	devices providers.DeviceProvider
}

func NewListDevices(devices providers.DeviceProvider) ListDevices {
	return &listDevicesImpl{devices: devices}
}

func (uc *listDevicesImpl) Execute(ctx context.Context, req ListDevicesRequest) ([]entities.Device, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	return uc.devices.ListDevices(ctx, req.OrgID, req.RampID)
}
