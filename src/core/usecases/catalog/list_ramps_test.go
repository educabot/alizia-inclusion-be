package catalog_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/educabot/alizia-inclusion-be/src/core/entities"
	"github.com/educabot/alizia-inclusion-be/src/core/providers"
	mockproviders "github.com/educabot/alizia-inclusion-be/src/mocks/providers"
	"github.com/educabot/alizia-inclusion-be/src/core/usecases/catalog"
	"github.com/educabot/alizia-inclusion-be/src/testutil"
)

func TestListRamps_ReturnsRampsForValidOrg(t *testing.T) {
	ctx := context.Background()
	expected := []entities.Ramp{
		testutil.NewRamp(1, "Ramp 1"),
		testutil.NewRamp(2, "Ramp 2"),
	}
	ramps := new(mockproviders.MockRampProvider)
	ramps.On("ListRamps", ctx, testutil.TestOrgID).Return(expected, nil)

	got, err := catalog.NewListRamps(ramps).Execute(ctx, catalog.ListRampsRequest{OrgID: testutil.TestOrgID})

	assert.NoError(t, err)
	assert.Equal(t, expected, got)
	ramps.AssertExpectations(t)
}

func TestListRamps_RejectsNilOrgID(t *testing.T) {
	ramps := new(mockproviders.MockRampProvider)

	_, err := catalog.NewListRamps(ramps).Execute(context.Background(), catalog.ListRampsRequest{OrgID: uuid.Nil})

	assert.ErrorIs(t, err, providers.ErrValidation)
	ramps.AssertNotCalled(t, "ListRamps", mock.Anything, mock.Anything)
}

func TestListRamps_PropagatesProviderError(t *testing.T) {
	ctx := context.Background()
	ramps := new(mockproviders.MockRampProvider)
	ramps.On("ListRamps", ctx, testutil.TestOrgID).Return(nil, errDB)

	_, err := catalog.NewListRamps(ramps).Execute(ctx, catalog.ListRampsRequest{OrgID: testutil.TestOrgID})

	assert.ErrorIs(t, err, errDB)
	ramps.AssertExpectations(t)
}
