package dashboard

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/educabot/alizia-inclusion-be/src/core/providers"
)

type GetMetricsRequest struct {
	OrgID uuid.UUID
}

func (r GetMetricsRequest) Validate() error {
	if r.OrgID == uuid.Nil {
		return errOrgIDRequired
	}
	return nil
}

type DeviceUsageResponse struct {
	DeviceID   int64  `json:"device_id"`
	DeviceName string `json:"device_name"`
	Count      int    `json:"count"`
}

type GetMetricsResponse struct {
	TotalStudents        int                   `json:"total_students"`
	StudentsWithProfiles int                   `json:"students_with_profiles"`
	TotalAdaptations     int                   `json:"total_adaptations"`
	AdaptationsByStatus  map[string]int        `json:"adaptations_by_status"`
	AdaptationsByType    map[string]int        `json:"adaptations_by_type"`
	TopUsedDevices       []DeviceUsageResponse `json:"top_used_devices"`
	AdaptationsThisWeek  int                   `json:"adaptations_this_week"`
	ClassroomCount       int                   `json:"classroom_count"`
}

type GetMetrics interface {
	Execute(ctx context.Context, req GetMetricsRequest) (*GetMetricsResponse, error)
}

type getMetricsImpl struct {
	students    providers.StudentProvider
	adaptations providers.AdaptationProvider
	classrooms  providers.ClassroomProvider
}

func NewGetMetrics(students providers.StudentProvider, adaptations providers.AdaptationProvider, classrooms providers.ClassroomProvider) GetMetrics {
	return &getMetricsImpl{students: students, adaptations: adaptations, classrooms: classrooms}
}

func (uc *getMetricsImpl) Execute(ctx context.Context, req GetMetricsRequest) (*GetMetricsResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	students, err := uc.students.List(ctx, req.OrgID)
	if err != nil {
		return nil, err
	}

	withProfiles := 0
	for i := range students {
		if students[i].Profile != nil {
			withProfiles++
		}
	}

	adaptations, err := uc.adaptations.List(ctx, req.OrgID, providers.AdaptationFilter{})
	if err != nil {
		return nil, err
	}

	byStatus := make(map[string]int)
	byType := make(map[string]int)
	for i := range adaptations {
		a := &adaptations[i]
		byStatus[a.Status]++
		if a.AdaptationType != "" {
			byType[a.AdaptationType]++
		}
	}

	classrooms, err := uc.classrooms.List(ctx, req.OrgID)
	if err != nil {
		return nil, err
	}

	weekAgo := time.Now().AddDate(0, 0, -7)
	thisWeek, _ := uc.adaptations.CountSince(ctx, req.OrgID, weekAgo)

	topDevicesRaw, _ := uc.adaptations.TopDevices(ctx, req.OrgID, 5)
	topDevices := make([]DeviceUsageResponse, len(topDevicesRaw))
	for i, d := range topDevicesRaw {
		topDevices[i] = DeviceUsageResponse{
			DeviceID:   d.DeviceID,
			DeviceName: d.DeviceName,
			Count:      d.Count,
		}
	}

	return &GetMetricsResponse{
		TotalStudents:        len(students),
		StudentsWithProfiles: withProfiles,
		TotalAdaptations:     len(adaptations),
		AdaptationsByStatus:  byStatus,
		AdaptationsByType:    byType,
		TopUsedDevices:       topDevices,
		AdaptationsThisWeek:  thisWeek,
		ClassroomCount:       len(classrooms),
	}, nil
}
