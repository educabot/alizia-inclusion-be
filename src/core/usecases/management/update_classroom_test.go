package management_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/educabot/alizia-inclusion-be/src/core/providers"
	mockproviders "github.com/educabot/alizia-inclusion-be/src/core/providers/mocks"
	"github.com/educabot/alizia-inclusion-be/src/core/usecases/management"
	"github.com/educabot/alizia-inclusion-be/src/testutil"
)

func TestUpdateClassroom_UpdatesClassroom(t *testing.T) {
	ctx := context.Background()
	newName := "Updated Name"
	newGrade := "4"
	newSection := "C"
	existing := testutil.NewClassroom(1, "Old Name")

	classrooms := new(mockproviders.MockClassroomProvider)
	classrooms.On("Get", ctx, testutil.TestOrgID, int64(1)).Return(&existing, nil)
	classrooms.On("Update", ctx, mock.AnythingOfType("*entities.Classroom")).Return(nil)

	uc := management.NewUpdateClassroom(classrooms)
	got, err := uc.Execute(ctx, management.UpdateClassroomRequest{
		OrgID:       testutil.TestOrgID,
		ClassroomID: 1,
		Name:        &newName,
		Grade:       &newGrade,
		Section:     &newSection,
	})

	assert.NoError(t, err)
	assert.NotNil(t, got)
	assert.Equal(t, newName, got.Name)
	assert.Equal(t, &newGrade, got.Grade)
	assert.Equal(t, &newSection, got.Section)
	classrooms.AssertExpectations(t)
}

func TestUpdateClassroom_RejectsNilOrgID(t *testing.T) {
	classrooms := new(mockproviders.MockClassroomProvider)
	uc := management.NewUpdateClassroom(classrooms)

	_, err := uc.Execute(context.Background(), management.UpdateClassroomRequest{
		OrgID:       uuid.Nil,
		ClassroomID: 1,
	})

	assert.ErrorIs(t, err, providers.ErrValidation)
	classrooms.AssertNotCalled(t, "Get", mock.Anything, mock.Anything, mock.Anything)
}

func TestUpdateClassroom_RejectsZeroClassroomID(t *testing.T) {
	classrooms := new(mockproviders.MockClassroomProvider)
	uc := management.NewUpdateClassroom(classrooms)

	_, err := uc.Execute(context.Background(), management.UpdateClassroomRequest{
		OrgID:       testutil.TestOrgID,
		ClassroomID: 0,
	})

	assert.ErrorIs(t, err, providers.ErrValidation)
	classrooms.AssertNotCalled(t, "Get", mock.Anything, mock.Anything, mock.Anything)
}

func TestUpdateClassroom_ReturnsNotFound(t *testing.T) {
	ctx := context.Background()
	classrooms := new(mockproviders.MockClassroomProvider)
	classrooms.On("Get", ctx, testutil.TestOrgID, int64(99)).Return(nil, errClassroomNotFound)

	uc := management.NewUpdateClassroom(classrooms)
	newName := "New Name"
	got, err := uc.Execute(ctx, management.UpdateClassroomRequest{
		OrgID:       testutil.TestOrgID,
		ClassroomID: 99,
		Name:        &newName,
	})

	assert.ErrorIs(t, err, providers.ErrNotFound)
	assert.Nil(t, got)
	classrooms.AssertNotCalled(t, "Update", mock.Anything, mock.Anything)
	classrooms.AssertExpectations(t)
}
