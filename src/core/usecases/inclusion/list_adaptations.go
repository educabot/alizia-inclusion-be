package inclusion

import (
	"context"

	"github.com/google/uuid"

	"github.com/educabot/alizia-inclusion-be/src/core/entities"
	"github.com/educabot/alizia-inclusion-be/src/core/providers"
)

type ListAdaptationsRequest struct {
	OrgID     uuid.UUID
	StudentID *int64
	// TeacherID hace privado el listado por docente (HU-4); DeviceID/Query filtran
	// por material de valija / tema.
	TeacherID *int64
	DeviceID  *int64
	Query     string
}

func (r ListAdaptationsRequest) Validate() error {
	if r.OrgID == uuid.Nil {
		return errOrgIDRequired
	}
	return nil
}

type ListAdaptations interface {
	Execute(ctx context.Context, req ListAdaptationsRequest) ([]entities.Adaptation, error)
}

type listAdaptationsImpl struct {
	adaptations providers.AdaptationProvider
}

func NewListAdaptations(adaptations providers.AdaptationProvider) ListAdaptations {
	return &listAdaptationsImpl{adaptations: adaptations}
}

func (uc *listAdaptationsImpl) Execute(ctx context.Context, req ListAdaptationsRequest) ([]entities.Adaptation, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	return uc.adaptations.List(ctx, req.OrgID, providers.AdaptationFilter{
		StudentID: req.StudentID,
		TeacherID: req.TeacherID,
		DeviceID:  req.DeviceID,
		Query:     req.Query,
	})
}
