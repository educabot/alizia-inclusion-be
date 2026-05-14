package catalog

import (
	"context"

	"github.com/google/uuid"

	"github.com/educabot/alizia-inclusion-be/src/core/entities"
	"github.com/educabot/alizia-inclusion-be/src/core/providers"
)

type ListRampsRequest struct {
	OrgID uuid.UUID
}

func (r ListRampsRequest) Validate() error {
	if r.OrgID == uuid.Nil {
		return errOrgIDRequired
	}
	return nil
}

type ListRamps interface {
	Execute(ctx context.Context, req ListRampsRequest) ([]entities.Ramp, error)
}

type listRampsImpl struct {
	ramps providers.RampProvider
}

func NewListRamps(ramps providers.RampProvider) ListRamps {
	return &listRampsImpl{ramps: ramps}
}

func (uc *listRampsImpl) Execute(ctx context.Context, req ListRampsRequest) ([]entities.Ramp, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	return uc.ramps.ListRamps(ctx, req.OrgID)
}
