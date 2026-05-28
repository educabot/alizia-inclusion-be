package inclusion_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/google/uuid"

	"github.com/educabot/alizia-inclusion-be/src/core/entities"
	"github.com/educabot/alizia-inclusion-be/src/core/providers"
	mockproviders "github.com/educabot/alizia-inclusion-be/src/core/providers/mocks"
	"github.com/educabot/alizia-inclusion-be/src/core/usecases/inclusion"
	"github.com/educabot/alizia-inclusion-be/src/testutil"
)

func TestListAdaptations_ReturnsAllAdaptations(t *testing.T) {
	ctx := context.Background()
	want := []entities.Adaptation{
		testutil.NewAdaptation(1, 1, 1),
		testutil.NewAdaptation(2, 2, 1),
	}
	adaptations := new(mockproviders.MockAdaptationProvider)
	adaptations.On("List", ctx, testutil.TestOrgID, (*int64)(nil)).Return(want, nil)

	got, err := inclusion.NewListAdaptations(adaptations).Execute(ctx, inclusion.ListAdaptationsRequest{
		OrgID:     testutil.TestOrgID,
		StudentID: nil,
	})

	require.NoError(t, err)
	assert.Equal(t, want, got)
	adaptations.AssertExpectations(t)
}

func TestListAdaptations_FiltersByStudent(t *testing.T) {
	ctx := context.Background()
	wantStudentID := int64(1)
	want := []entities.Adaptation{
		testutil.NewAdaptation(1, wantStudentID, 1),
	}
	adaptations := new(mockproviders.MockAdaptationProvider)
	adaptations.On("List", ctx, testutil.TestOrgID, testutil.Ptr(wantStudentID)).Return(want, nil)

	got, err := inclusion.NewListAdaptations(adaptations).Execute(ctx, inclusion.ListAdaptationsRequest{
		OrgID:     testutil.TestOrgID,
		StudentID: testutil.Ptr(wantStudentID),
	})

	require.NoError(t, err)
	assert.Len(t, got, 1)
	adaptations.AssertExpectations(t)
}

func TestListAdaptations_RejectsNilOrgID(t *testing.T) {
	ctx := context.Background()
	adaptations := new(mockproviders.MockAdaptationProvider)

	_, err := inclusion.NewListAdaptations(adaptations).Execute(ctx, inclusion.ListAdaptationsRequest{
		OrgID: uuid.Nil,
	})

	assert.ErrorIs(t, err, providers.ErrValidation)
	adaptations.AssertNotCalled(t, "List", mock.Anything, mock.Anything, mock.Anything)
}
