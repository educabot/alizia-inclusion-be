package catalog_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/educabot/alizia-inclusion-be/src/core/entities"
	"github.com/educabot/alizia-inclusion-be/src/core/providers"
	mockproviders "github.com/educabot/alizia-inclusion-be/src/core/providers/mocks"
	"github.com/educabot/alizia-inclusion-be/src/core/usecases/catalog"
	"github.com/educabot/alizia-inclusion-be/src/testutil"
)

func TestListDevices_ReturnsAllDevices(t *testing.T) {
	ctx := context.Background()
	expected := []entities.Device{
		testutil.NewDevice(1, 1, "Device 1"),
		testutil.NewDevice(2, 1, "Device 2"),
	}
	devices := new(mockproviders.MockDeviceProvider)
	devices.On("ListDevices", ctx, testutil.TestOrgID, (*int64)(nil)).Return(expected, nil)

	got, err := catalog.NewListDevices(devices).Execute(ctx, catalog.ListDevicesRequest{OrgID: testutil.TestOrgID, RampID: nil})

	assert.NoError(t, err)
	assert.Equal(t, expected, got)
	devices.AssertExpectations(t)
}

func TestListDevices_FiltersByRampID(t *testing.T) {
	ctx := context.Background()
	wantRampID := int64(1)
	expected := []entities.Device{testutil.NewDevice(1, wantRampID, "Device 1")}
	devices := new(mockproviders.MockDeviceProvider)
	devices.On("ListDevices", ctx, testutil.TestOrgID, testutil.Ptr(wantRampID)).Return(expected, nil)

	got, err := catalog.NewListDevices(devices).Execute(ctx, catalog.ListDevicesRequest{
		OrgID:  testutil.TestOrgID,
		RampID: testutil.Ptr(wantRampID),
	})

	assert.NoError(t, err)
	assert.Len(t, got, 1)
	assert.Equal(t, expected[0].ID, got[0].ID)
	devices.AssertExpectations(t)
}

func TestListDevices_RejectsNilOrgID(t *testing.T) {
	devices := new(mockproviders.MockDeviceProvider)

	_, err := catalog.NewListDevices(devices).Execute(context.Background(), catalog.ListDevicesRequest{OrgID: uuid.Nil, RampID: nil})

	assert.ErrorIs(t, err, providers.ErrValidation)
	devices.AssertNotCalled(t, "ListDevices", mock.Anything, mock.Anything, mock.Anything)
}
