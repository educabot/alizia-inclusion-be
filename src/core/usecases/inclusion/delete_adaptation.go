package inclusion

import (
	"context"

	"github.com/google/uuid"

	"github.com/educabot/alizia-inclusion-be/src/core/providers"
)

type DeleteAdaptationRequest struct {
	OrgID        uuid.UUID
	AdaptationID int64
}

func (r DeleteAdaptationRequest) Validate() error {
	if r.OrgID == uuid.Nil {
		return errOrgIDRequired
	}
	if r.AdaptationID <= 0 {
		return errAdaptationIDRequired
	}
	return nil
}

type DeleteAdaptation interface {
	Execute(ctx context.Context, req DeleteAdaptationRequest) error
}

type deleteAdaptationImpl struct {
	adaptations providers.AdaptationProvider
}

func NewDeleteAdaptation(adaptations providers.AdaptationProvider) DeleteAdaptation {
	return &deleteAdaptationImpl{adaptations: adaptations}
}

func (uc *deleteAdaptationImpl) Execute(ctx context.Context, req DeleteAdaptationRequest) error {
	if err := req.Validate(); err != nil {
		return err
	}
	return uc.adaptations.Delete(ctx, req.OrgID, req.AdaptationID)
}
