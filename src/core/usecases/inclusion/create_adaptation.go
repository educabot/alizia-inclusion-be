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
	Title               string
	Subject             string
	ActivityDescription *string
	AdaptationStrategy  *string
	AdaptationType      string
	Notes               *string
	// Origen IA (HU-4): de qué conversación/mensaje salió y si el docente la editó.
	SourceConversationID *int64
	SourceMessageID      *int64
	WasEdited            bool
}

const defaultAdaptationType = "actividad_adaptada"

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

	adaptationType := req.AdaptationType
	if adaptationType == "" {
		adaptationType = defaultAdaptationType
	}

	adaptation := &entities.Adaptation{
		OrganizationID:       req.OrgID,
		StudentID:            req.StudentID,
		TeacherID:            req.TeacherID,
		DeviceID:             req.DeviceID,
		Title:                req.Title,
		Subject:              req.Subject,
		ActivityDescription:  req.ActivityDescription,
		AdaptationStrategy:   req.AdaptationStrategy,
		AdaptationType:       adaptationType,
		Notes:                req.Notes,
		Status:               "en_curso",
		SourceConversationID: req.SourceConversationID,
		SourceMessageID:      req.SourceMessageID,
		WasEdited:            req.WasEdited,
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
