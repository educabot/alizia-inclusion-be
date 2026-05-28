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

func TestGetClassroom_ReturnsClassroom(t *testing.T) {
	ctx := context.Background()
	want := testutil.NewClassroom(1, "1A")
	classrooms := new(mockproviders.MockClassroomProvider)
	classrooms.On("Get", ctx, testutil.TestOrgID, int64(1)).Return(&want, nil)

	uc := management.NewGetClassroom(classrooms)
	got, err := uc.Execute(ctx, management.GetClassroomRequest{
		OrgID:       testutil.TestOrgID,
		ClassroomID: 1,
	})

	assert.NoError(t, err)
	assert.Equal(t, &want, got)
	classrooms.AssertExpectations(t)
}

func TestGetClassroom_RejectsNilOrgID(t *testing.T) {
	classrooms := new(mockproviders.MockClassroomProvider)
	uc := management.NewGetClassroom(classrooms)

	_, err := uc.Execute(context.Background(), management.GetClassroomRequest{
		OrgID:       uuid.Nil,
		ClassroomID: 1,
	})

	assert.ErrorIs(t, err, providers.ErrValidation)
	classrooms.AssertNotCalled(t, "Get", mock.Anything, mock.Anything, mock.Anything)
}

func TestGetClassroom_RejectsZeroClassroomID(t *testing.T) {
	classrooms := new(mockproviders.MockClassroomProvider)
	uc := management.NewGetClassroom(classrooms)

	_, err := uc.Execute(context.Background(), management.GetClassroomRequest{
		OrgID:       testutil.TestOrgID,
		ClassroomID: 0,
	})

	assert.ErrorIs(t, err, providers.ErrValidation)
	classrooms.AssertNotCalled(t, "Get", mock.Anything, mock.Anything, mock.Anything)
}

func TestGetClassroom_ReturnsNotFound(t *testing.T) {
	ctx := context.Background()
	classrooms := new(mockproviders.MockClassroomProvider)
	classrooms.On("Get", ctx, testutil.TestOrgID, int64(99)).Return(nil, errClassroomNotFound)

	uc := management.NewGetClassroom(classrooms)
	got, err := uc.Execute(ctx, management.GetClassroomRequest{
		OrgID:       testutil.TestOrgID,
		ClassroomID: 99,
	})

	assert.ErrorIs(t, err, providers.ErrNotFound)
	assert.Nil(t, got)
	classrooms.AssertExpectations(t)
}
