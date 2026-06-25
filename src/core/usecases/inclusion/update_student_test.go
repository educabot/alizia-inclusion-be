package inclusion_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/educabot/alizia-inclusion-be/src/core/providers"
	mockproviders "github.com/educabot/alizia-inclusion-be/src/mocks/providers"
	"github.com/educabot/alizia-inclusion-be/src/core/usecases/inclusion"
	"github.com/educabot/alizia-inclusion-be/src/testutil"
)

func TestUpdateStudent_UpdatesStudent(t *testing.T) {
	students := new(mockproviders.MockStudentProvider)
	ctx := context.Background()
	existing := testutil.NewStudent(1, 1, "Old Name")
	students.On("GetStudent", ctx, testutil.TestOrgID, int64(1)).Return(&existing, nil)
	students.On("Update", ctx, mock.AnythingOfType("*entities.Student")).Return(nil)

	newName := "New Name"
	got, err := inclusion.NewUpdateStudent(students).Execute(ctx, inclusion.UpdateStudentRequest{
		OrgID:     testutil.TestOrgID,
		StudentID: 1,
		Name:      &newName,
	})

	require.NoError(t, err)
	assert.Equal(t, newName, got.Name)
	students.AssertExpectations(t)
}

func TestUpdateStudent_RejectsNilOrgID(t *testing.T) {
	students := new(mockproviders.MockStudentProvider)

	_, err := inclusion.NewUpdateStudent(students).Execute(context.Background(), inclusion.UpdateStudentRequest{
		OrgID: uuid.Nil, StudentID: 1,
	})

	assert.ErrorIs(t, err, providers.ErrValidation)
	students.AssertNotCalled(t, "GetStudent", mock.Anything, mock.Anything, mock.Anything)
	students.AssertNotCalled(t, "Update", mock.Anything, mock.Anything)
}

func TestUpdateStudent_RejectsZeroStudentID(t *testing.T) {
	students := new(mockproviders.MockStudentProvider)

	_, err := inclusion.NewUpdateStudent(students).Execute(context.Background(), inclusion.UpdateStudentRequest{
		OrgID: testutil.TestOrgID, StudentID: 0,
	})

	assert.ErrorIs(t, err, providers.ErrValidation)
	students.AssertNotCalled(t, "GetStudent", mock.Anything, mock.Anything, mock.Anything)
	students.AssertNotCalled(t, "Update", mock.Anything, mock.Anything)
}

func TestUpdateStudent_ReturnsNotFound(t *testing.T) {
	students := new(mockproviders.MockStudentProvider)
	ctx := context.Background()
	students.On("GetStudent", ctx, testutil.TestOrgID, int64(999)).Return(nil, errStudentNotFound)

	_, err := inclusion.NewUpdateStudent(students).Execute(ctx, inclusion.UpdateStudentRequest{
		OrgID: testutil.TestOrgID, StudentID: 999,
	})

	assert.ErrorIs(t, err, providers.ErrNotFound)
	students.AssertExpectations(t)
	students.AssertNotCalled(t, "Update", mock.Anything, mock.Anything)
}
