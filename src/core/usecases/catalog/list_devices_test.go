package catalog_test

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"

	"github.com/educabot/alizia-inclusion-be/src/core/entities"
	"github.com/educabot/alizia-inclusion-be/src/core/providers"
	"github.com/educabot/alizia-inclusion-be/src/core/providers/mocks"
	"github.com/educabot/alizia-inclusion-be/src/core/usecases/catalog"
	"github.com/educabot/alizia-inclusion-be/src/testutil"
)

func TestListDevices(t *testing.T) {
	ctx := context.Background()

	t.Run("returns all devices", func(t *testing.T) {
		expected := []entities.Device{
			testutil.NewDevice(1, 1, "Device 1"),
			testutil.NewDevice(2, 1, "Device 2"),
		}

		mock := &mocks.MockDeviceProvider{
			ListDevicesFn: func(_ context.Context, _ uuid.UUID, _ *int64) ([]entities.Device, error) {
				return expected, nil
			},
		}

		req := catalog.ListDevicesRequest{OrgID: testutil.TestOrgID, RampID: nil}
		got, err := catalog.NewListDevices(mock).Execute(ctx, req)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if len(got) != len(expected) {
			t.Errorf("got %d devices, want %d", len(got), len(expected))
		}
		for i, d := range got {
			if d.ID != expected[i].ID || d.Name != expected[i].Name {
				t.Errorf("device[%d] = {ID:%d Name:%q}, want {ID:%d Name:%q}",
					i, d.ID, d.Name, expected[i].ID, expected[i].Name)
			}
		}
	})

	t.Run("filters by ramp_id", func(t *testing.T) {
		wantRampID := int64(1)
		var capturedRampID *int64

		mock := &mocks.MockDeviceProvider{
			ListDevicesFn: func(_ context.Context, _ uuid.UUID, rampID *int64) ([]entities.Device, error) {
				capturedRampID = rampID
				return []entities.Device{testutil.NewDevice(1, wantRampID, "Device 1")}, nil
			},
		}

		req := catalog.ListDevicesRequest{
			OrgID:  testutil.TestOrgID,
			RampID: testutil.Ptr(wantRampID),
		}
		got, err := catalog.NewListDevices(mock).Execute(ctx, req)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if len(got) != 1 {
			t.Errorf("got %d devices, want 1", len(got))
		}
		if capturedRampID == nil {
			t.Fatal("rampID was not passed to mock")
		}
		if *capturedRampID != wantRampID {
			t.Errorf("mock received rampID %d, want %d", *capturedRampID, wantRampID)
		}
	})

	t.Run("rejects nil org_id", func(t *testing.T) {
		called := false
		mock := &mocks.MockDeviceProvider{
			ListDevicesFn: func(_ context.Context, _ uuid.UUID, _ *int64) ([]entities.Device, error) {
				called = true
				return nil, nil
			},
		}

		req := catalog.ListDevicesRequest{OrgID: uuid.Nil, RampID: nil}
		_, err := catalog.NewListDevices(mock).Execute(ctx, req)
		if err == nil {
			t.Error("expected validation error, got nil")
		}
		if !errors.Is(err, providers.ErrValidation) {
			t.Errorf("expected ErrValidation, got: %v", err)
		}
		if called {
			t.Error("mock should not have been called for invalid request")
		}
	})
}
