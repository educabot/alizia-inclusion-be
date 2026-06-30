package inclusion

import (
	"context"

	"github.com/google/uuid"

	"github.com/educabot/alizia-inclusion-be/src/core/entities"
	"github.com/educabot/alizia-inclusion-be/src/core/providers"
)

type CreateAdaptationRequest struct {
	OrgID               uuid.UUID
	StudentID           *int64
	TeacherID           int64
	DeviceID            *int64
	DeviceIDs           []int64
	Title               string
	Subject             string
	ActivityDescription *string
	AdaptationStrategy  *string
	AdaptationType      string
	Outcome             *string
	Notes               *string
	Steps               entities.AdaptationSteps
	RampID              *int64
	// Origen opcional cuando la adaptación se guarda desde el chat (GAP B).
	SourceConversationID *int64
	SourceMessageID      *int64
}

const defaultAdaptationType = "actividad_adaptada"

// validAdaptationTypes es el enum permitido para adaptation_type. Vacío → default;
// un valor no vacío debe pertenecer al set (si no, error de validación).
var validAdaptationTypes = map[string]struct{}{
	"actividad_adaptada":  {},
	"material_nuevo":      {},
	"estrategia_aula":     {},
	"situacion_emergente": {},
}

func (r CreateAdaptationRequest) Validate() error {
	if r.OrgID == uuid.Nil {
		return errOrgIDRequired
	}
	if r.StudentID != nil && *r.StudentID <= 0 {
		return errStudentIDRequired
	}
	if r.TeacherID <= 0 {
		return errTeacherIDRequired
	}
	// Subject (materia) es opcional: el diseño del flujo del docente descarta "materia"
	// (se usa el curso/grado del alumno). El guardado desde el chat tampoco lo envía.
	if r.AdaptationType != "" {
		if _, ok := validAdaptationTypes[r.AdaptationType]; !ok {
			return errInvalidAdaptationType
		}
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
		Outcome:              req.Outcome,
		Notes:                req.Notes,
		Steps:                req.Steps,
		RampID:               req.RampID,
		Status:               "en_curso",
		SourceConversationID: req.SourceConversationID,
		SourceMessageID:      req.SourceMessageID,
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
