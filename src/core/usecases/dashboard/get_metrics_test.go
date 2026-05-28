package dashboard_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"

	"github.com/educabot/alizia-inclusion-be/src/core/entities"
	"github.com/educabot/alizia-inclusion-be/src/core/providers"
	"github.com/educabot/alizia-inclusion-be/src/core/providers/mocks"
	"github.com/educabot/alizia-inclusion-be/src/core/usecases/dashboard"
	"github.com/educabot/alizia-inclusion-be/src/testutil"
)

func TestGetMetrics_ReturnsMetrics(t *testing.T) {
	ctx := context.Background()
	students := []entities.Student{
		testutil.NewStudentWithProfile(1, 10, "Maria", []string{"lectura"}),
		testutil.NewStudentWithProfile(2, 10, "Pedro", []string{"atencion"}),
		testutil.NewStudent(3, 10, "Luis"),
	}
	adaptations := []entities.Adaptation{
		testutil.NewAdaptation(1, 1, 5),
		func() entities.Adaptation {
			a := testutil.NewAdaptation(2, 2, 5)
			a.Status = "finalizada"
			a.AdaptationType = "material_nuevo"
			return a
		}(),
	}
	topDevices := []providers.DeviceUsageStat{
		{DeviceID: 1, DeviceName: "Timer", Count: 3},
	}
	classrooms := []entities.Classroom{
		testutil.NewClassroom(10, "1A"),
		testutil.NewClassroom(11, "2B"),
	}

	studentMock := &mocks.MockStudentProvider{
		ListFn: func(_ context.Context, _ uuid.UUID) ([]entities.Student, error) {
			return students, nil
		},
	}
	adaptationMock := &mocks.MockAdaptationProvider{
		ListFn: func(_ context.Context, _ uuid.UUID, _ *int64) ([]entities.Adaptation, error) {
			return adaptations, nil
		},
		CountSinceFn: func(_ context.Context, _ uuid.UUID, _ time.Time) (int, error) {
			return 5, nil
		},
		TopDevicesFn: func(_ context.Context, _ uuid.UUID, _ int) ([]providers.DeviceUsageStat, error) {
			return topDevices, nil
		},
	}
	classroomMock := &mocks.MockClassroomProvider{
		ListFn: func(_ context.Context, _ uuid.UUID) ([]entities.Classroom, error) {
			return classrooms, nil
		},
	}
	uc := dashboard.NewGetMetrics(studentMock, adaptationMock, classroomMock)

	got, err := uc.Execute(ctx, dashboard.GetMetricsRequest{OrgID: testutil.TestOrgID})

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if got == nil {
		t.Fatal("expected non-nil response, got nil")
	}
	if got.TotalStudents != 3 {
		t.Errorf("TotalStudents = %d, want 3", got.TotalStudents)
	}
	if got.StudentsWithProfiles != 2 {
		t.Errorf("StudentsWithProfiles = %d, want 2", got.StudentsWithProfiles)
	}
	if got.TotalAdaptations != 2 {
		t.Errorf("TotalAdaptations = %d, want 2", got.TotalAdaptations)
	}
	if got.AdaptationsThisWeek != 5 {
		t.Errorf("AdaptationsThisWeek = %d, want 5", got.AdaptationsThisWeek)
	}
	if got.ClassroomCount != 2 {
		t.Errorf("ClassroomCount = %d, want 2", got.ClassroomCount)
	}
	if len(got.TopUsedDevices) != 1 {
		t.Fatalf("TopUsedDevices len = %d, want 1", len(got.TopUsedDevices))
	}
	if got.TopUsedDevices[0].DeviceID != 1 {
		t.Errorf("TopUsedDevices[0].DeviceID = %d, want 1", got.TopUsedDevices[0].DeviceID)
	}
	if got.TopUsedDevices[0].DeviceName != "Timer" {
		t.Errorf("TopUsedDevices[0].DeviceName = %q, want Timer", got.TopUsedDevices[0].DeviceName)
	}
	if got.TopUsedDevices[0].Count != 3 {
		t.Errorf("TopUsedDevices[0].Count = %d, want 3", got.TopUsedDevices[0].Count)
	}
}

func TestGetMetrics_RejectsNilOrgID(t *testing.T) {
	ctx := context.Background()
	uc := dashboard.NewGetMetrics(
		&mocks.MockStudentProvider{},
		&mocks.MockAdaptationProvider{},
		&mocks.MockClassroomProvider{},
	)

	_, err := uc.Execute(ctx, dashboard.GetMetricsRequest{OrgID: uuid.Nil})

	if err == nil {
		t.Fatal("expected validation error, got nil")
	}
	if !errors.Is(err, providers.ErrValidation) {
		t.Errorf("expected ErrValidation, got %v", err)
	}
}
