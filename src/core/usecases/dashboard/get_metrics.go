package dashboard

import (
	"context"

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

type GetMetricsResponse struct {
	TotalStudents        int            `json:"total_students"`
	StudentsWithProfiles int            `json:"students_with_profiles"`
	TotalAdaptations     int            `json:"total_adaptations"`
	AdaptationsByStatus  map[string]int `json:"adaptations_by_status"`
	ClassroomCount       int            `json:"classroom_count"`
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
	for _, s := range students {
		if s.Profile != nil {
			withProfiles++
		}
	}

	adaptations, err := uc.adaptations.List(ctx, req.OrgID, nil)
	if err != nil {
		return nil, err
	}

	byStatus := make(map[string]int)
	for _, a := range adaptations {
		byStatus[a.Status]++
	}

	classrooms, err := uc.classrooms.List(ctx, req.OrgID)
	if err != nil {
		return nil, err
	}

	return &GetMetricsResponse{
		TotalStudents:        len(students),
		StudentsWithProfiles: withProfiles,
		TotalAdaptations:     len(adaptations),
		AdaptationsByStatus:  byStatus,
		ClassroomCount:       len(classrooms),
	}, nil
}
