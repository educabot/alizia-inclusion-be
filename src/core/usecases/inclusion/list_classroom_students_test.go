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

func TestListClassroomStudents_ReturnsStudents(t *testing.T) {
	students := new(mockproviders.MockStudentProvider)
	ctx := context.Background()
	expected := []entities.Student{
		testutil.NewStudent(1, 5, "Ana"),
		testutil.NewStudent(2, 5, "Lucas"),
	}
	students.On("ListByClassroom", ctx, testutil.TestOrgID, int64(5)).Return(expected, nil)

	got, err := inclusion.NewListClassroomStudents(students).Execute(ctx, inclusion.ListClassroomStudentsRequest{
		OrgID:       testutil.TestOrgID,
		ClassroomID: 5,
	})

	require.NoError(t, err)
	assert.Len(t, got, 2)
	students.AssertExpectations(t)
}

func TestListClassroomStudents_RejectsNilOrgID(t *testing.T) {
	students := new(mockproviders.MockStudentProvider)

	_, err := inclusion.NewListClassroomStudents(students).Execute(context.Background(), inclusion.ListClassroomStudentsRequest{
		OrgID: uuid.Nil, ClassroomID: 1,
	})

	assert.ErrorIs(t, err, providers.ErrValidation)
	students.AssertNotCalled(t, "ListByClassroom", mock.Anything, mock.Anything, mock.Anything)
}

func TestListClassroomStudents_RejectsZeroClassroomID(t *testing.T) {
	students := new(mockproviders.MockStudentProvider)

	_, err := inclusion.NewListClassroomStudents(students).Execute(context.Background(), inclusion.ListClassroomStudentsRequest{
		OrgID: testutil.TestOrgID, ClassroomID: 0,
	})

	assert.ErrorIs(t, err, providers.ErrValidation)
	students.AssertNotCalled(t, "ListByClassroom", mock.Anything, mock.Anything, mock.Anything)
}
