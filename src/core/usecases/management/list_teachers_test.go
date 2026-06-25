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

func TestListTeachers_ReturnsTeachers(t *testing.T) {
	ctx := context.Background()
	want := []entities.User{
		testutil.NewUser(1, "Ana Garcia"),
		testutil.NewUser(2, "Luis Perez"),
	}
	users := new(mockproviders.MockUserProvider)
	users.On("ListByRole", ctx, testutil.TestOrgID, "teacher").Return(want, nil)

	uc := management.NewListTeachers(users)
	got, err := uc.Execute(ctx, management.ListTeachersRequest{OrgID: testutil.TestOrgID})

	assert.NoError(t, err)
	assert.Len(t, got, 2)
	assert.Equal(t, want[0].ID, got[0].ID)
	users.AssertExpectations(t)
}

func TestListTeachers_RejectsNilOrgID(t *testing.T) {
	users := new(mockproviders.MockUserProvider)
	uc := management.NewListTeachers(users)

	_, err := uc.Execute(context.Background(), management.ListTeachersRequest{OrgID: uuid.Nil})

	assert.ErrorIs(t, err, providers.ErrValidation)
	users.AssertNotCalled(t, "ListByRole", mock.Anything, mock.Anything, mock.Anything)
}
