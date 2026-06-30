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
	StudentID           *int64
	DeviceID            *int64
	DeviceIDs           *[]int64
	Title               *string
	Subject             *string
	ActivityDescription *string
	AdaptationStrategy  *string
	AdaptationType      *string
	Outcome             *string
	Notes               *string
	Status              *string
	Steps               *entities.AdaptationSteps
	RampID              *int64
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

	if req.StudentID != nil {
		existing.StudentID = req.StudentID
	}
	if req.DeviceID != nil {
		existing.DeviceID = req.DeviceID
	}
	if req.Title != nil {
		existing.Title = *req.Title
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
	if req.AdaptationType != nil {
		existing.AdaptationType = *req.AdaptationType
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
	if req.Steps != nil {
		existing.Steps = *req.Steps
	}
	if req.RampID != nil {
		existing.RampID = req.RampID
	}

	if err := uc.adaptations.Update(ctx, existing); err != nil {
		return nil, err
	}

	if req.DeviceIDs != nil {
		if err := uc.adaptations.SetDevices(ctx, existing.ID, *req.DeviceIDs); err != nil {
			return nil, err
		}
	}

	refreshed, err := uc.adaptations.Get(ctx, req.OrgID, req.AdaptationID)
	if err != nil {
		return nil, err
	}
	return refreshed, nil
}
