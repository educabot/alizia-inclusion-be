package catalog

import (
	"context"

	"github.com/google/uuid"

	"github.com/educabot/alizia-inclusion-be/src/core/entities"
	"github.com/educabot/alizia-inclusion-be/src/core/providers"
)

type GetRampRequest struct {
	OrgID  uuid.UUID
	RampID int64
}

func (r GetRampRequest) Validate() error {
	if r.OrgID == uuid.Nil {
		return errOrgIDRequired
	}
	if r.RampID <= 0 {
		return errRampIDRequired
	}
	return nil
}

type GetRamp interface {
	Execute(ctx context.Context, req GetRampRequest) (*entities.Ramp, error)
}

type getRampImpl struct {
	ramps providers.RampProvider
}

func NewGetRamp(ramps providers.RampProvider) GetRamp {
	return &getRampImpl{ramps: ramps}
}

func (uc *getRampImpl) Execute(ctx context.Context, req GetRampRequest) (*entities.Ramp, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	return uc.ramps.GetRamp(ctx, req.OrgID, req.RampID)
}
