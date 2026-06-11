package catalog_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/educabot/alizia-inclusion-be/src/core/providers"
	mockproviders "github.com/educabot/alizia-inclusion-be/src/mocks/providers"
	"github.com/educabot/alizia-inclusion-be/src/core/usecases/catalog"
	"github.com/educabot/alizia-inclusion-be/src/testutil"
)

func TestGetRamp_ReturnsRamp(t *testing.T) {
	ctx := context.Background()
	expected := testutil.NewRamp(1, "Ramp 1")
	ramps := new(mockproviders.MockRampProvider)
	ramps.On("GetRamp", ctx, testutil.TestOrgID, int64(1)).Return(&expected, nil)

	got, err := catalog.NewGetRamp(ramps).Execute(ctx, catalog.GetRampRequest{OrgID: testutil.TestOrgID, RampID: 1})

	assert.NoError(t, err)
	assert.Equal(t, &expected, got)
	ramps.AssertExpectations(t)
}

func TestGetRamp_RejectsNilOrgID(t *testing.T) {
	ramps := new(mockproviders.MockRampProvider)

	_, err := catalog.NewGetRamp(ramps).Execute(context.Background(), catalog.GetRampRequest{OrgID: uuid.Nil, RampID: 1})

	assert.ErrorIs(t, err, providers.ErrValidation)
	ramps.AssertNotCalled(t, "GetRamp", mock.Anything, mock.Anything, mock.Anything)
}

func TestGetRamp_RejectsZeroRampID(t *testing.T) {
	ramps := new(mockproviders.MockRampProvider)

	_, err := catalog.NewGetRamp(ramps).Execute(context.Background(), catalog.GetRampRequest{OrgID: testutil.TestOrgID, RampID: 0})

	assert.ErrorIs(t, err, providers.ErrValidation)
	ramps.AssertNotCalled(t, "GetRamp", mock.Anything, mock.Anything, mock.Anything)
}

func TestGetRamp_ReturnsNotFound(t *testing.T) {
	ctx := context.Background()
	ramps := new(mockproviders.MockRampProvider)
	ramps.On("GetRamp", ctx, testutil.TestOrgID, int64(99)).Return(nil, errRampNotFound)

	got, err := catalog.NewGetRamp(ramps).Execute(ctx, catalog.GetRampRequest{OrgID: testutil.TestOrgID, RampID: 99})

	assert.ErrorIs(t, err, providers.ErrNotFound)
	assert.Nil(t, got)
	ramps.AssertExpectations(t)
}
