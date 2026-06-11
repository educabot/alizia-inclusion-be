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

func TestGetAdaptation_ReturnsAdaptation(t *testing.T) {
	adaptations := new(mockproviders.MockAdaptationProvider)
	ctx := context.Background()
	expected := testutil.NewAdaptation(1, 1, 1)
	adaptations.On("Get", ctx, testutil.TestOrgID, int64(1)).Return(&expected, nil)

	got, err := inclusion.NewGetAdaptation(adaptations).Execute(ctx, inclusion.GetAdaptationRequest{
		OrgID:        testutil.TestOrgID,
		AdaptationID: 1,
	})

	require.NoError(t, err)
	assert.Equal(t, expected.ID, got.ID)
	adaptations.AssertExpectations(t)
}

func TestGetAdaptation_RejectsNilOrgID(t *testing.T) {
	adaptations := new(mockproviders.MockAdaptationProvider)

	_, err := inclusion.NewGetAdaptation(adaptations).Execute(context.Background(), inclusion.GetAdaptationRequest{
		OrgID:        uuid.Nil,
		AdaptationID: 1,
	})

	assert.ErrorIs(t, err, providers.ErrValidation)
	adaptations.AssertNotCalled(t, "Get", mock.Anything, mock.Anything, mock.Anything)
}

func TestGetAdaptation_RejectsZeroAdaptationID(t *testing.T) {
	adaptations := new(mockproviders.MockAdaptationProvider)

	_, err := inclusion.NewGetAdaptation(adaptations).Execute(context.Background(), inclusion.GetAdaptationRequest{
		OrgID:        testutil.TestOrgID,
		AdaptationID: 0,
	})

	assert.ErrorIs(t, err, providers.ErrValidation)
	adaptations.AssertNotCalled(t, "Get", mock.Anything, mock.Anything, mock.Anything)
}

func TestGetAdaptation_ReturnsNotFound(t *testing.T) {
	adaptations := new(mockproviders.MockAdaptationProvider)
	ctx := context.Background()
	adaptations.On("Get", ctx, testutil.TestOrgID, int64(99)).Return(nil, errAdaptationNotFound)

	_, err := inclusion.NewGetAdaptation(adaptations).Execute(ctx, inclusion.GetAdaptationRequest{
		OrgID:        testutil.TestOrgID,
		AdaptationID: 99,
	})

	assert.ErrorIs(t, err, providers.ErrNotFound)
	adaptations.AssertExpectations(t)
}
