package management_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/educabot/alizia-inclusion-be/src/core/entities"
	"github.com/educabot/alizia-inclusion-be/src/core/providers"
	mockproviders "github.com/educabot/alizia-inclusion-be/src/mocks/providers"
	"github.com/educabot/alizia-inclusion-be/src/core/usecases/management"
	"github.com/educabot/alizia-inclusion-be/src/testutil"
)

func TestListClassrooms_ReturnsClassrooms(t *testing.T) {
	ctx := context.Background()
	want := []entities.Classroom{
		testutil.NewClassroom(1, "1A"),
		testutil.NewClassroom(2, "2B"),
	}
	classrooms := new(mockproviders.MockClassroomProvider)
	classrooms.On("List", ctx, testutil.TestOrgID).Return(want, nil)

	uc := management.NewListClassrooms(classrooms)
	got, err := uc.Execute(ctx, management.ListClassroomsRequest{OrgID: testutil.TestOrgID})

	assert.NoError(t, err)
	assert.Len(t, got, 2)
	assert.Equal(t, want[0].ID, got[0].ID)
	classrooms.AssertExpectations(t)
}

func TestListClassrooms_RejectsNilOrgID(t *testing.T) {
	classrooms := new(mockproviders.MockClassroomProvider)
	uc := management.NewListClassrooms(classrooms)

	_, err := uc.Execute(context.Background(), management.ListClassroomsRequest{OrgID: uuid.Nil})

	assert.ErrorIs(t, err, providers.ErrValidation)
	classrooms.AssertNotCalled(t, "List", mock.Anything, mock.Anything)
}
