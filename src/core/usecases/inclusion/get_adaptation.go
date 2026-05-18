package inclusion

import (
	"context"

	"github.com/google/uuid"

	"github.com/educabot/alizia-inclusion-be/src/core/entities"
	"github.com/educabot/alizia-inclusion-be/src/core/providers"
)

type GetAdaptationRequest struct {
	OrgID        uuid.UUID
	AdaptationID int64
}

func (r GetAdaptationRequest) Validate() error {
	if r.OrgID == uuid.Nil {
		return errOrgIDRequired
	}
	if r.AdaptationID <= 0 {
		return errAdaptationIDRequired
	}
	return nil
}

type GetAdaptation interface {
	Execute(ctx context.Context, req GetAdaptationRequest) (*entities.Adaptation, error)
}

type getAdaptationImpl struct {
	adaptations providers.AdaptationProvider
}

func NewGetAdaptation(adaptations providers.AdaptationProvider) GetAdaptation {
	return &getAdaptationImpl{adaptations: adaptations}
}

func (uc *getAdaptationImpl) Execute(ctx context.Context, req GetAdaptationRequest) (*entities.Adaptation, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	return uc.adaptations.Get(ctx, req.OrgID, req.AdaptationID)
}
