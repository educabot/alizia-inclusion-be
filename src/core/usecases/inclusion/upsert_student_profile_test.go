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

func TestUpsertStudentProfile_UpsertsProfile(t *testing.T) {
	students := new(mockproviders.MockStudentProvider)
	profiles := new(mockproviders.MockStudentProfileProvider)
	ctx := context.Background()
	existing := testutil.NewStudent(1, 1, "Lucas")
	students.On("GetStudent", ctx, testutil.TestOrgID, int64(1)).Return(&existing, nil)
	profiles.On("Upsert", ctx, mock.AnythingOfType("*entities.StudentProfile")).
		Run(func(args mock.Arguments) {
			p, ok := args.Get(1).(*entities.StudentProfile)
			require.True(t, ok)
			p.ID = 1
		}).
		Return(nil)
	profiles.On("GetByStudentID", ctx, int64(1)).Return(&entities.StudentProfile{
		ID:           1,
		StudentID:    1,
		IsTransitory: true,
		Difficulties: []string{"distraccion", "motricidad_fina"},
	}, nil)

	got, err := inclusion.NewUpsertStudentProfile(students, profiles).Execute(ctx, inclusion.UpsertStudentProfileRequest{
		OrgID:        testutil.TestOrgID,
		StudentID:    1,
		IsTransitory: true,
		Difficulties: []string{"distraccion", "motricidad_fina"},
	})

	require.NoError(t, err)
	require.NotNil(t, got)
	assert.True(t, got.IsTransitory)
	assert.Len(t, got.Difficulties, 2)
	students.AssertExpectations(t)
	profiles.AssertExpectations(t)
}

func TestUpsertStudentProfile_RejectsNilOrgID(t *testing.T) {
	students := new(mockproviders.MockStudentProvider)
	profiles := new(mockproviders.MockStudentProfileProvider)

	_, err := inclusion.NewUpsertStudentProfile(students, profiles).Execute(context.Background(), inclusion.UpsertStudentProfileRequest{
		OrgID: uuid.Nil, StudentID: 1, Difficulties: []string{},
	})

	assert.ErrorIs(t, err, providers.ErrValidation)
	students.AssertNotCalled(t, "GetStudent", mock.Anything, mock.Anything, mock.Anything)
	profiles.AssertNotCalled(t, "Upsert", mock.Anything, mock.Anything)
}

func TestUpsertStudentProfile_RejectsZeroStudentID(t *testing.T) {
	students := new(mockproviders.MockStudentProvider)
	profiles := new(mockproviders.MockStudentProfileProvider)

	_, err := inclusion.NewUpsertStudentProfile(students, profiles).Execute(context.Background(), inclusion.UpsertStudentProfileRequest{
		OrgID: testutil.TestOrgID, StudentID: 0, Difficulties: []string{},
	})

	assert.ErrorIs(t, err, providers.ErrValidation)
	students.AssertNotCalled(t, "GetStudent", mock.Anything, mock.Anything, mock.Anything)
	profiles.AssertNotCalled(t, "Upsert", mock.Anything, mock.Anything)
}

func TestUpsertStudentProfile_ReturnsErrorIfStudentNotFound(t *testing.T) {
	students := new(mockproviders.MockStudentProvider)
	profiles := new(mockproviders.MockStudentProfileProvider)
	ctx := context.Background()
	students.On("GetStudent", ctx, testutil.TestOrgID, int64(999)).Return(nil, errStudentNotFound)

	_, err := inclusion.NewUpsertStudentProfile(students, profiles).Execute(ctx, inclusion.UpsertStudentProfileRequest{
		OrgID: testutil.TestOrgID, StudentID: 999, Difficulties: []string{},
	})

	assert.ErrorIs(t, err, providers.ErrNotFound)
	students.AssertExpectations(t)
	profiles.AssertNotCalled(t, "Upsert", mock.Anything, mock.Anything)
}
