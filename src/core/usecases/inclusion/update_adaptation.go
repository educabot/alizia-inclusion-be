package inclusion

import (
	"context"

	"github.com/google/uuid"

	"github.com/educabot/alizia-inclusion-be/src/core/entities"
	"github.com/educabot/alizia-inclusion-be/src/core/providers"
)

type UpdateAdaptationRequest struct {
	OrgID               uuid.UUID
	AdaptationID        int64
	DeviceID            *int64
	Subject             *string
	ActivityDescription *string
	AdaptationStrategy  *string
	Outcome             *string
	Notes               *string
	Status              *string
}

var validStatuses = map[string]bool{
	"en_curso":     true,
	"probado":      true,
	"funciono":     true,
	"para_ajustar": true,
}

func (r UpdateAdaptationRequest) Validate() error {
	if r.OrgID == uuid.Nil {
		return errOrgIDRequired
	}
	if r.AdaptationID <= 0 {
		return errAdaptationIDRequired
	}
	if r.Status != nil && !validStatuses[*r.Status] {
		return errInvalidStatus
	}
	return nil
}

type UpdateAdaptation interface {
	Execute(ctx context.Context, req UpdateAdaptationRequest) (*entities.Adaptation, error)
}

type updateAdaptationImpl struct {
	adaptations providers.AdaptationProvider
}

func NewUpdateAdaptation(adaptations providers.AdaptationProvider) UpdateAdaptation {
	return &updateAdaptationImpl{adaptations: adaptations}
}

func (uc *updateAdaptationImpl) Execute(ctx context.Context, req UpdateAdaptationRequest) (*entities.Adaptation, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	existing, err := uc.adaptations.Get(ctx, req.OrgID, req.AdaptationID)
	if err != nil {
		return nil, err
	}

	if req.DeviceID != nil {
		existing.DeviceID = req.DeviceID
	}
	if req.Subject != nil {
		existing.Subject = *req.Subject
	}
	if req.ActivityDescription != nil {
		existing.ActivityDescription = req.ActivityDescription
	}
	if req.AdaptationStrategy != nil {
		existing.AdaptationStrategy = req.AdaptationStrategy
	}
	if req.Outcome != nil {
		existing.Outcome = req.Outcome
	}
	if req.Notes != nil {
		existing.Notes = req.Notes
	}
	if req.Status != nil {
		existing.Status = *req.Status
	}

	if err := uc.adaptations.Update(ctx, existing); err != nil {
		return nil, err
	}
	return existing, nil
}
