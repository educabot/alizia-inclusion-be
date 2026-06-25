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

func TestDeleteStudent_DeletesStudent(t *testing.T) {
	students := new(mockproviders.MockStudentProvider)
	ctx := context.Background()
	students.On("Delete", ctx, testutil.TestOrgID, int64(5)).Return(nil)

	err := inclusion.NewDeleteStudent(students).Execute(ctx, inclusion.DeleteStudentRequest{
		OrgID:     testutil.TestOrgID,
		StudentID: 5,
	})

	require.NoError(t, err)
	students.AssertExpectations(t)
}

func TestDeleteStudent_RejectsNilOrgID(t *testing.T) {
	students := new(mockproviders.MockStudentProvider)

	err := inclusion.NewDeleteStudent(students).Execute(context.Background(), inclusion.DeleteStudentRequest{
		OrgID:     uuid.Nil,
		StudentID: 1,
	})

	assert.ErrorIs(t, err, providers.ErrValidation)
	students.AssertNotCalled(t, "Delete", mock.Anything, mock.Anything, mock.Anything)
}

func TestDeleteStudent_RejectsZeroStudentID(t *testing.T) {
	students := new(mockproviders.MockStudentProvider)

	err := inclusion.NewDeleteStudent(students).Execute(context.Background(), inclusion.DeleteStudentRequest{
		OrgID:     testutil.TestOrgID,
		StudentID: 0,
	})

	assert.ErrorIs(t, err, providers.ErrValidation)
	students.AssertNotCalled(t, "Delete", mock.Anything, mock.Anything, mock.Anything)
}
