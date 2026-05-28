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

func TestListStudents_ReturnsAllStudents(t *testing.T) {
	students := new(mockproviders.MockStudentProvider)
	ctx := context.Background()
	expected := []entities.Student{
		testutil.NewStudent(1, 1, "Lucas"),
		testutil.NewStudent(2, 1, "Ana"),
	}
	students.On("List", ctx, testutil.TestOrgID).Return(expected, nil)

	got, err := inclusion.NewListStudents(students).Execute(ctx, inclusion.ListStudentsRequest{
		OrgID: testutil.TestOrgID,
	})

	require.NoError(t, err)
	assert.Len(t, got, 2)
	students.AssertExpectations(t)
}

func TestListStudents_FiltersByClassroom(t *testing.T) {
	students := new(mockproviders.MockStudentProvider)
	ctx := context.Background()
	classroomID := int64(5)
	students.On("ListByClassroom", ctx, testutil.TestOrgID, classroomID).
		Return([]entities.Student{testutil.NewStudent(1, classroomID, "Lucas")}, nil)

	got, err := inclusion.NewListStudents(students).Execute(ctx, inclusion.ListStudentsRequest{
		OrgID:       testutil.TestOrgID,
		ClassroomID: &classroomID,
	})

	require.NoError(t, err)
	assert.Len(t, got, 1)
	students.AssertExpectations(t)
}

func TestListStudents_RejectsNilOrgID(t *testing.T) {
	students := new(mockproviders.MockStudentProvider)

	_, err := inclusion.NewListStudents(students).Execute(context.Background(), inclusion.ListStudentsRequest{
		OrgID: uuid.Nil,
	})

	assert.ErrorIs(t, err, providers.ErrValidation)
	students.AssertNotCalled(t, "List", mock.Anything, mock.Anything)
	students.AssertNotCalled(t, "ListByClassroom", mock.Anything, mock.Anything, mock.Anything)
}
