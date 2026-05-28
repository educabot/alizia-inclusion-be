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

func TestUpdateAdaptation_UpdatesAdaptation(t *testing.T) {
	adaptations := new(mockproviders.MockAdaptationProvider)
	ctx := context.Background()
	existing := testutil.NewAdaptation(1, 1, 1)
	updated := testutil.NewAdaptation(1, 1, 1)
	updated.Subject = "Lengua"
	updated.Status = "probado"

	adaptations.On("Get", ctx, testutil.TestOrgID, int64(1)).Return(&existing, nil).Once()
	adaptations.On("Update", ctx, mock.AnythingOfType("*entities.Adaptation")).Return(nil)
	adaptations.On("Get", ctx, testutil.TestOrgID, int64(1)).Return(&updated, nil).Once()

	got, err := inclusion.NewUpdateAdaptation(adaptations).Execute(ctx, inclusion.UpdateAdaptationRequest{
		OrgID:        testutil.TestOrgID,
		AdaptationID: 1,
		Subject:      testutil.Ptr("Lengua"),
		Status:       testutil.Ptr("probado"),
	})

	require.NoError(t, err)
	require.NotNil(t, got)
	assert.Equal(t, "Lengua", got.Subject)
	assert.Equal(t, "probado", got.Status)
	adaptations.AssertExpectations(t)
}

func TestUpdateAdaptation_AppliesTitleWhenProvided(t *testing.T) {
	adaptations := new(mockproviders.MockAdaptationProvider)
	ctx := context.Background()
	existing := testutil.NewAdaptation(1, 1, 1)
	var captured *entities.Adaptation
	adaptations.On("Get", ctx, testutil.TestOrgID, int64(1)).Return(&existing, nil)
	adaptations.On("Update", ctx, mock.AnythingOfType("*entities.Adaptation")).
		Run(func(args mock.Arguments) {
			captured = args.Get(1).(*entities.Adaptation)
		}).
		Return(nil)

	_, err := inclusion.NewUpdateAdaptation(adaptations).Execute(ctx, inclusion.UpdateAdaptationRequest{
		OrgID:        testutil.TestOrgID,
		AdaptationID: 1,
		Title:        testutil.Ptr("Secuencia con apoyos visuales"),
	})

	require.NoError(t, err)
	require.NotNil(t, captured)
	assert.Equal(t, "Secuencia con apoyos visuales", captured.Title)
}

func TestUpdateAdaptation_RejectsNilOrgID(t *testing.T) {
	adaptations := new(mockproviders.MockAdaptationProvider)

	_, err := inclusion.NewUpdateAdaptation(adaptations).Execute(context.Background(), inclusion.UpdateAdaptationRequest{
		OrgID:        uuid.Nil,
		AdaptationID: 1,
	})

	assert.ErrorIs(t, err, providers.ErrValidation)
	adaptations.AssertNotCalled(t, "Get", mock.Anything, mock.Anything, mock.Anything)
	adaptations.AssertNotCalled(t, "Update", mock.Anything, mock.Anything)
}

func TestUpdateAdaptation_RejectsZeroAdaptationID(t *testing.T) {
	adaptations := new(mockproviders.MockAdaptationProvider)

	_, err := inclusion.NewUpdateAdaptation(adaptations).Execute(context.Background(), inclusion.UpdateAdaptationRequest{
		OrgID:        testutil.TestOrgID,
		AdaptationID: 0,
	})

	assert.ErrorIs(t, err, providers.ErrValidation)
	adaptations.AssertNotCalled(t, "Get", mock.Anything, mock.Anything, mock.Anything)
}

func TestUpdateAdaptation_RejectsInvalidStatus(t *testing.T) {
	adaptations := new(mockproviders.MockAdaptationProvider)

	_, err := inclusion.NewUpdateAdaptation(adaptations).Execute(context.Background(), inclusion.UpdateAdaptationRequest{
		OrgID:        testutil.TestOrgID,
		AdaptationID: 1,
		Status:       testutil.Ptr("invalid"),
	})

	assert.ErrorIs(t, err, providers.ErrValidation)
	adaptations.AssertNotCalled(t, "Get", mock.Anything, mock.Anything, mock.Anything)
}

func TestUpdateAdaptation_ReturnsNotFound(t *testing.T) {
	adaptations := new(mockproviders.MockAdaptationProvider)
	ctx := context.Background()
	adaptations.On("Get", ctx, testutil.TestOrgID, int64(99)).Return(nil, errAdaptationNotFound)

	_, err := inclusion.NewUpdateAdaptation(adaptations).Execute(ctx, inclusion.UpdateAdaptationRequest{
		OrgID:        testutil.TestOrgID,
		AdaptationID: 99,
	})

	assert.ErrorIs(t, err, providers.ErrNotFound)
	adaptations.AssertExpectations(t)
	adaptations.AssertNotCalled(t, "Update", mock.Anything, mock.Anything)
}
