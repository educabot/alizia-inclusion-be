package management_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/educabot/alizia-inclusion-be/src/core/providers"
	mockproviders "github.com/educabot/alizia-inclusion-be/src/mocks/providers"
	"github.com/educabot/alizia-inclusion-be/src/core/usecases/management"
	"github.com/educabot/alizia-inclusion-be/src/testutil"
)

func TestCreateClassroom_CreatesClassroom(t *testing.T) {
	ctx := context.Background()
	grade := "3"
	section := "B"

	classrooms := new(mockproviders.MockClassroomProvider)
	classrooms.On("Create", ctx, mock.AnythingOfType("*entities.Classroom")).Return(nil)

	uc := management.NewCreateClassroom(classrooms)
	got, err := uc.Execute(ctx, management.CreateClassroomRequest{
		OrgID:   testutil.TestOrgID,
		Name:    "3B Math",
		Grade:   &grade,
		Section: &section,
	})

	assert.NoError(t, err)
	assert.NotNil(t, got)
	assert.Equal(t, testutil.TestOrgID, got.OrganizationID)
	assert.Equal(t, "3B Math", got.Name)
	assert.Equal(t, &grade, got.Grade)
	assert.Equal(t, &section, got.Section)
	classrooms.AssertExpectations(t)
}

func TestCreateClassroom_RejectsNilOrgID(t *testing.T) {
	classrooms := new(mockproviders.MockClassroomProvider)
	uc := management.NewCreateClassroom(classrooms)

	_, err := uc.Execute(context.Background(), management.CreateClassroomRequest{
		OrgID: uuid.Nil,
		Name:  "1A",
	})

	assert.ErrorIs(t, err, providers.ErrValidation)
	classrooms.AssertNotCalled(t, "Create", mock.Anything, mock.Anything)
}

func TestCreateClassroom_RejectsEmptyName(t *testing.T) {
	classrooms := new(mockproviders.MockClassroomProvider)
	uc := management.NewCreateClassroom(classrooms)

	_, err := uc.Execute(context.Background(), management.CreateClassroomRequest{
		OrgID: testutil.TestOrgID,
		Name:  "",
	})

	assert.ErrorIs(t, err, providers.ErrValidation)
	classrooms.AssertNotCalled(t, "Create", mock.Anything, mock.Anything)
}

func TestCreateClassroom_PropagatesProviderError(t *testing.T) {
	ctx := context.Background()
	classrooms := new(mockproviders.MockClassroomProvider)
	classrooms.On("Create", ctx, mock.AnythingOfType("*entities.Classroom")).Return(errDB)

	uc := management.NewCreateClassroom(classrooms)
	_, err := uc.Execute(ctx, management.CreateClassroomRequest{
		OrgID: testutil.TestOrgID,
		Name:  "1A",
	})

	assert.ErrorIs(t, err, errDB)
	classrooms.AssertExpectations(t)
}
