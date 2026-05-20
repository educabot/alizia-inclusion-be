package inclusion

import (
	"context"

	"github.com/google/uuid"

	"github.com/educabot/alizia-inclusion-be/src/core/entities"
	"github.com/educabot/alizia-inclusion-be/src/core/providers"
)

type CreateAdaptationRequest struct {
	OrgID               uuid.UUID
	StudentID           int64
	TeacherID           int64
	DeviceID            *int64
	DeviceIDs           []int64
	Subject             string
	ActivityDescription *string
	AdaptationStrategy  *string
	AdaptationType      string
	Notes               *string
}

func (r CreateAdaptationRequest) Validate() error {
	if r.OrgID == uuid.Nil {
		return errOrgIDRequired
	}
	if r.StudentID <= 0 {
		return errStudentIDRequired
	}
	if r.TeacherID <= 0 {
		return errTeacherIDRequired
	}
	if r.Subject == "" {
		return errSubjectRequired
	}
	return nil
}

type CreateAdaptation interface {
	Execute(ctx context.Context, req CreateAdaptationRequest) (*entities.Adaptation, error)
}

type createAdaptationImpl struct {
	adaptations providers.AdaptationProvider
}

func NewCreateAdaptation(adaptations providers.AdaptationProvider) CreateAdaptation {
	return &createAdaptationImpl{adaptations: adaptations}
}

func (uc *createAdaptationImpl) Execute(ctx context.Context, req CreateAdaptationRequest) (*entities.Adaptation, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	adaptation := &entities.Adaptation{
		OrganizationID:      req.OrgID,
		StudentID:           req.StudentID,
		TeacherID:           req.TeacherID,
		DeviceID:            req.DeviceID,
		Subject:             req.Subject,
		ActivityDescription: req.ActivityDescription,
		AdaptationStrategy:  req.AdaptationStrategy,
		AdaptationType:      req.AdaptationType,
		Notes:               req.Notes,
		Status:              "en_curso",
	}

	if err := uc.adaptations.Create(ctx, adaptation); err != nil {
		return nil, err
	}

	if len(req.DeviceIDs) > 0 {
		if err := uc.adaptations.SetDevices(ctx, adaptation.ID, req.DeviceIDs); err != nil {
			return nil, err
		}
	}

	refreshed, err := uc.adaptations.Get(ctx, req.OrgID, adaptation.ID)
	if err != nil {
		return nil, err
	}
	return refreshed, nil
}
