package inclusion_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/educabot/alizia-inclusion-be/src/core/providers"
	mockproviders "github.com/educabot/alizia-inclusion-be/src/core/providers/mocks"
	"github.com/educabot/alizia-inclusion-be/src/core/usecases/inclusion"
	"github.com/educabot/alizia-inclusion-be/src/testutil"
)

func TestDeleteAdaptation_DeletesAdaptation(t *testing.T) {
	adaptations := new(mockproviders.MockAdaptationProvider)
	ctx := context.Background()
	adaptations.On("Delete", ctx, testutil.TestOrgID, int64(1)).Return(nil)

	err := inclusion.NewDeleteAdaptation(adaptations).Execute(ctx, inclusion.DeleteAdaptationRequest{
		OrgID:        testutil.TestOrgID,
		AdaptationID: 1,
	})

	require.NoError(t, err)
	adaptations.AssertExpectations(t)
}

func TestDeleteAdaptation_RejectsNilOrgID(t *testing.T) {
	adaptations := new(mockproviders.MockAdaptationProvider)

	err := inclusion.NewDeleteAdaptation(adaptations).Execute(context.Background(), inclusion.DeleteAdaptationRequest{
		OrgID:        uuid.Nil,
		AdaptationID: 1,
	})

	assert.ErrorIs(t, err, providers.ErrValidation)
	adaptations.AssertNotCalled(t, "Delete", mock.Anything, mock.Anything, mock.Anything)
}

func TestDeleteAdaptation_RejectsZeroAdaptationID(t *testing.T) {
	adaptations := new(mockproviders.MockAdaptationProvider)

	err := inclusion.NewDeleteAdaptation(adaptations).Execute(context.Background(), inclusion.DeleteAdaptationRequest{
		OrgID:        testutil.TestOrgID,
		AdaptationID: 0,
	})

	assert.ErrorIs(t, err, providers.ErrValidation)
	adaptations.AssertNotCalled(t, "Delete", mock.Anything, mock.Anything, mock.Anything)
}
