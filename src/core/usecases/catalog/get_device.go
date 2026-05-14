package catalog

import (
	"context"

	"github.com/google/uuid"

	"github.com/educabot/alizia-inclusion-be/src/core/entities"
	"github.com/educabot/alizia-inclusion-be/src/core/providers"
)

type GetDeviceRequest struct {
	OrgID    uuid.UUID
	DeviceID int64
}

func (r GetDeviceRequest) Validate() error {
	if r.OrgID == uuid.Nil {
		return errOrgIDRequired
	}
	if r.DeviceID <= 0 {
		return errDeviceIDRequired
	}
	return nil
}

type GetDevice interface {
	Execute(ctx context.Context, req GetDeviceRequest) (*entities.Device, error)
}

type getDeviceImpl struct {
	devices providers.DeviceProvider
}

func NewGetDevice(devices providers.DeviceProvider) GetDevice {
	return &getDeviceImpl{devices: devices}
}

func (uc *getDeviceImpl) Execute(ctx context.Context, req GetDeviceRequest) (*entities.Device, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	return uc.devices.GetDevice(ctx, req.OrgID, req.DeviceID)
}
