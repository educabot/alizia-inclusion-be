package providers

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/educabot/alizia-inclusion-be/src/core/entities"
)

type DeviceUsageStat struct {
	DeviceID   int64
	DeviceName string
	Count      int
}

// AdaptationFilter acota el listado de recursos (HU-4). TeacherID hace privado
// el recurso por docente; StudentID/DeviceID/Query filtran por alumno / material
// de valija / tema (texto en subject/title/actividad).
type AdaptationFilter struct {
	StudentID *int64
	TeacherID *int64
	DeviceID  *int64
	Query     string
}

type AdaptationProvider interface {
	List(ctx context.Context, orgID uuid.UUID, filter AdaptationFilter) ([]entities.Adaptation, error)
	Get(ctx context.Context, orgID uuid.UUID, id int64) (*entities.Adaptation, error)
	Create(ctx context.Context, adaptation *entities.Adaptation) error
	Update(ctx context.Context, adaptation *entities.Adaptation) error
	Delete(ctx context.Context, orgID uuid.UUID, id int64) error
	SetDevices(ctx context.Context, adaptationID int64, deviceIDs []int64) error
	CountSince(ctx context.Context, orgID uuid.UUID, since time.Time) (int, error)
	TopDevices(ctx context.Context, orgID uuid.UUID, limit int) ([]DeviceUsageStat, error)
}

type AdaptationResourceProvider interface {
	ListByAdaptation(ctx context.Context, adaptationID int64) ([]entities.AdaptationResource, error)
}
