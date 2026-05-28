package inclusion_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/educabot/alizia-inclusion-be/src/core/entities"
	"github.com/educabot/alizia-inclusion-be/src/core/providers"
	mockproviders "github.com/educabot/alizia-inclusion-be/src/core/providers/mocks"
	"github.com/educabot/alizia-inclusion-be/src/core/usecases/inclusion"
	"github.com/educabot/alizia-inclusion-be/src/testutil"
)

func TestCreateStudent_CreatesStudent(t *testing.T) {
	students := new(mockproviders.MockStudentProvider)
	ctx := context.Background()

	var captured *entities.Student
	students.On("Create", ctx, mock.AnythingOfType("*entities.Student")).
		Return(nil).
		Run(func(args mock.Arguments) {
			s, ok := args.Get(1).(*entities.Student)
			require.True(t, ok)
			s.ID = 10
			captured = s
		})

	got, err := inclusion.NewCreateStudent(students).Execute(ctx, inclusion.CreateStudentRequest{
		OrgID:       testutil.TestOrgID,
		ClassroomID: 1,
		Name:        "Juan",
	})

	require.NoError(t, err)
	require.NotNil(t, got)
	assert.Equal(t, "Juan", captured.Name)
	assert.Equal(t, int64(1), captured.ClassroomID)
	assert.Equal(t, int64(10), got.ID)
	students.AssertExpectations(t)
}

func TestCreateStudent_RejectsNilOrgID(t *testing.T) {
	students := new(mockproviders.MockStudentProvider)

	_, err := inclusion.NewCreateStudent(students).Execute(context.Background(), inclusion.CreateStudentRequest{
		OrgID: uuid.Nil, ClassroomID: 1, Name: "X",
	})

	assert.ErrorIs(t, err, providers.ErrValidation)
	students.AssertNotCalled(t, "Create", mock.Anything, mock.Anything)
}

func TestCreateStudent_RejectsZeroClassroomID(t *testing.T) {
	students := new(mockproviders.MockStudentProvider)

	_, err := inclusion.NewCreateStudent(students).Execute(context.Background(), inclusion.CreateStudentRequest{
		OrgID: testutil.TestOrgID, ClassroomID: 0, Name: "X",
	})

	assert.ErrorIs(t, err, providers.ErrValidation)
	students.AssertNotCalled(t, "Create", mock.Anything, mock.Anything)
}

func TestCreateStudent_RejectsEmptyName(t *testing.T) {
	students := new(mockproviders.MockStudentProvider)

	_, err := inclusion.NewCreateStudent(students).Execute(context.Background(), inclusion.CreateStudentRequest{
		OrgID: testutil.TestOrgID, ClassroomID: 1, Name: "",
	})

	assert.ErrorIs(t, err, providers.ErrValidation)
	students.AssertNotCalled(t, "Create", mock.Anything, mock.Anything)
}
