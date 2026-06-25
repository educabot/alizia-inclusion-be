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

func TestGetStudentProfile_ReturnsStudentWithProfile(t *testing.T) {
	students := new(mockproviders.MockStudentProvider)
	ctx := context.Background()
	expected := testutil.NewStudentWithProfile(1, 1, "Lucas", []string{"distraccion"})
	students.On("GetStudent", ctx, testutil.TestOrgID, int64(1)).Return(&expected, nil)

	got, err := inclusion.NewGetStudentProfile(students).Execute(ctx, inclusion.GetStudentProfileRequest{
		OrgID:     testutil.TestOrgID,
		StudentID: 1,
	})

	require.NoError(t, err)
	require.NotNil(t, got.Profile)
	assert.Equal(t, "distraccion", got.Profile.Difficulties[0])
	students.AssertExpectations(t)
}

func TestGetStudentProfile_RejectsNilOrgID(t *testing.T) {
	students := new(mockproviders.MockStudentProvider)

	_, err := inclusion.NewGetStudentProfile(students).Execute(context.Background(), inclusion.GetStudentProfileRequest{
		OrgID: uuid.Nil, StudentID: 1,
	})

	assert.ErrorIs(t, err, providers.ErrValidation)
	students.AssertNotCalled(t, "GetStudent", mock.Anything, mock.Anything, mock.Anything)
}

func TestGetStudentProfile_RejectsZeroStudentID(t *testing.T) {
	students := new(mockproviders.MockStudentProvider)

	_, err := inclusion.NewGetStudentProfile(students).Execute(context.Background(), inclusion.GetStudentProfileRequest{
		OrgID: testutil.TestOrgID, StudentID: 0,
	})

	assert.ErrorIs(t, err, providers.ErrValidation)
	students.AssertNotCalled(t, "GetStudent", mock.Anything, mock.Anything, mock.Anything)
}

func TestGetStudentProfile_ReturnsNotFound(t *testing.T) {
	students := new(mockproviders.MockStudentProvider)
	ctx := context.Background()
	students.On("GetStudent", ctx, testutil.TestOrgID, int64(999)).Return(nil, errStudentNotFound)

	_, err := inclusion.NewGetStudentProfile(students).Execute(ctx, inclusion.GetStudentProfileRequest{
		OrgID: testutil.TestOrgID, StudentID: 999,
	})

	assert.ErrorIs(t, err, providers.ErrNotFound)
	students.AssertExpectations(t)
}
