package dashboard_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/educabot/alizia-inclusion-be/src/core/entities"
	"github.com/educabot/alizia-inclusion-be/src/core/providers"
	mockproviders "github.com/educabot/alizia-inclusion-be/src/core/providers/mocks"
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

	studentMock := new(mockproviders.MockStudentProvider)
	studentMock.On("List", ctx, testutil.TestOrgID).Return(students, nil)

	adaptationMock := new(mockproviders.MockAdaptationProvider)
	adaptationMock.On("List", ctx, testutil.TestOrgID, (*int64)(nil)).Return(adaptations, nil)
	adaptationMock.On("CountSince", ctx, testutil.TestOrgID, mock.AnythingOfType("time.Time")).Return(5, nil)
	adaptationMock.On("TopDevices", ctx, testutil.TestOrgID, 5).Return(topDevices, nil)

	classroomMock := new(mockproviders.MockClassroomProvider)
	classroomMock.On("List", ctx, testutil.TestOrgID).Return(classrooms, nil)

	uc := dashboard.NewGetMetrics(studentMock, adaptationMock, classroomMock)

	got, err := uc.Execute(ctx, dashboard.GetMetricsRequest{OrgID: testutil.TestOrgID})

	assert.NoError(t, err)
	assert.NotNil(t, got)
	assert.Equal(t, 3, got.TotalStudents)
	assert.Equal(t, 2, got.StudentsWithProfiles)
	assert.Equal(t, 2, got.TotalAdaptations)
	assert.Equal(t, 5, got.AdaptationsThisWeek)
	assert.Equal(t, 2, got.ClassroomCount)
	assert.Len(t, got.TopUsedDevices, 1)
	assert.Equal(t, int64(1), got.TopUsedDevices[0].DeviceID)
	assert.Equal(t, "Timer", got.TopUsedDevices[0].DeviceName)
	assert.Equal(t, 3, got.TopUsedDevices[0].Count)
	studentMock.AssertExpectations(t)
	adaptationMock.AssertExpectations(t)
	classroomMock.AssertExpectations(t)
}

func TestGetMetrics_RejectsNilOrgID(t *testing.T) {
	studentMock := new(mockproviders.MockStudentProvider)
	adaptationMock := new(mockproviders.MockAdaptationProvider)
	classroomMock := new(mockproviders.MockClassroomProvider)
	uc := dashboard.NewGetMetrics(studentMock, adaptationMock, classroomMock)

	_, err := uc.Execute(context.Background(), dashboard.GetMetricsRequest{OrgID: uuid.Nil})

	assert.ErrorIs(t, err, providers.ErrValidation)
	studentMock.AssertNotCalled(t, "List", mock.Anything, mock.Anything)
	adaptationMock.AssertNotCalled(t, "List", mock.Anything, mock.Anything, mock.Anything)
	classroomMock.AssertNotCalled(t, "List", mock.Anything, mock.Anything)
}
