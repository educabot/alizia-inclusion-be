package inclusion_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/educabot/alizia-inclusion-be/src/core/entities"
	"github.com/educabot/alizia-inclusion-be/src/core/providers"
	mockproviders "github.com/educabot/alizia-inclusion-be/src/core/providers/mocks"
	"github.com/educabot/alizia-inclusion-be/src/core/usecases/inclusion"
	"github.com/educabot/alizia-inclusion-be/src/testutil"
)

func TestListAdaptationResources_ReturnsResources(t *testing.T) {
	ctx := context.Background()
	want := []entities.AdaptationResource{
		testutil.NewAdaptationResource(1, 10),
		testutil.NewAdaptationResource(2, 10),
	}
	resources := new(mockproviders.MockAdaptationResourceProvider)
	resources.On("ListByAdaptation", ctx, int64(10)).Return(want, nil)

	got, err := inclusion.NewListAdaptationResources(resources).Execute(ctx, inclusion.ListAdaptationResourcesRequest{
		AdaptationID: 10,
	})

	require.NoError(t, err)
	assert.Len(t, got, 2)
	resources.AssertExpectations(t)
}

func TestListAdaptationResources_RejectsZeroAdaptationID(t *testing.T) {
	ctx := context.Background()
	resources := new(mockproviders.MockAdaptationResourceProvider)

	_, err := inclusion.NewListAdaptationResources(resources).Execute(ctx, inclusion.ListAdaptationResourcesRequest{
		AdaptationID: 0,
	})

	assert.ErrorIs(t, err, providers.ErrValidation)
	resources.AssertNotCalled(t, "ListByAdaptation")
}
