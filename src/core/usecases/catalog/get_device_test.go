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

func TestGetDevice_ReturnsDevice(t *testing.T) {
	ctx := context.Background()
	expected := testutil.NewDevice(1, 1, "Timer Visual")
	devices := new(mockproviders.MockDeviceProvider)
	devices.On("GetDevice", ctx, testutil.TestOrgID, int64(1)).Return(&expected, nil)

	got, err := catalog.NewGetDevice(devices).Execute(ctx, catalog.GetDeviceRequest{
		OrgID:    testutil.TestOrgID,
		DeviceID: 1,
	})

	assert.NoError(t, err)
	assert.Equal(t, &expected, got)
	devices.AssertExpectations(t)
}

func TestGetDevice_RejectsNilOrgID(t *testing.T) {
	devices := new(mockproviders.MockDeviceProvider)

	_, err := catalog.NewGetDevice(devices).Execute(context.Background(), catalog.GetDeviceRequest{
		OrgID:    uuid.Nil,
		DeviceID: 1,
	})

	assert.ErrorIs(t, err, providers.ErrValidation)
	devices.AssertNotCalled(t, "GetDevice", mock.Anything, mock.Anything, mock.Anything)
}

func TestGetDevice_RejectsZeroDeviceID(t *testing.T) {
	devices := new(mockproviders.MockDeviceProvider)

	_, err := catalog.NewGetDevice(devices).Execute(context.Background(), catalog.GetDeviceRequest{
		OrgID:    testutil.TestOrgID,
		DeviceID: 0,
	})

	assert.ErrorIs(t, err, providers.ErrValidation)
	devices.AssertNotCalled(t, "GetDevice", mock.Anything, mock.Anything, mock.Anything)
}

func TestGetDevice_ReturnsNotFound(t *testing.T) {
	ctx := context.Background()
	devices := new(mockproviders.MockDeviceProvider)
	devices.On("GetDevice", ctx, testutil.TestOrgID, int64(999)).Return(nil, errDevNotFound)

	got, err := catalog.NewGetDevice(devices).Execute(ctx, catalog.GetDeviceRequest{
		OrgID:    testutil.TestOrgID,
		DeviceID: 999,
	})

	assert.ErrorIs(t, err, providers.ErrNotFound)
	assert.Nil(t, got)
	devices.AssertExpectations(t)
}
