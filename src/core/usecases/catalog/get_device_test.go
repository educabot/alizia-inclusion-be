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

func TestGetDevice_ReturnsDevice(t *testing.T) {
	ctx := context.Background()
	expected := testutil.NewDevice(1, 1, "Timer Visual")
	mock := &mocks.MockDeviceProvider{
		GetDeviceFn: func(_ context.Context, _ uuid.UUID, id int64) (*entities.Device, error) {
			d := expected
			return &d, nil
		},
	}

	got, err := catalog.NewGetDevice(mock).Execute(ctx, catalog.GetDeviceRequest{
		OrgID:    testutil.TestOrgID,
		DeviceID: 1,
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.ID != expected.ID || got.Name != expected.Name {
		t.Errorf("got {ID:%d Name:%q}, want {ID:%d Name:%q}", got.ID, got.Name, expected.ID, expected.Name)
	}
}

func TestGetDevice_RejectsNilOrgID(t *testing.T) {
	ctx := context.Background()
	mock := &mocks.MockDeviceProvider{}

	_, err := catalog.NewGetDevice(mock).Execute(ctx, catalog.GetDeviceRequest{
		OrgID:    uuid.Nil,
		DeviceID: 1,
	})

	if !errors.Is(err, providers.ErrValidation) {
		t.Errorf("expected ErrValidation, got: %v", err)
	}
}

func TestGetDevice_RejectsZeroDeviceID(t *testing.T) {
	ctx := context.Background()
	mock := &mocks.MockDeviceProvider{}

	_, err := catalog.NewGetDevice(mock).Execute(ctx, catalog.GetDeviceRequest{
		OrgID:    testutil.TestOrgID,
		DeviceID: 0,
	})

	if !errors.Is(err, providers.ErrValidation) {
		t.Errorf("expected ErrValidation, got: %v", err)
	}
}

func TestGetDevice_ReturnsNotFound(t *testing.T) {
	ctx := context.Background()
	mock := &mocks.MockDeviceProvider{
		GetDeviceFn: func(_ context.Context, _ uuid.UUID, _ int64) (*entities.Device, error) {
			return nil, errDevNotFound
		},
	}

	_, err := catalog.NewGetDevice(mock).Execute(ctx, catalog.GetDeviceRequest{
		OrgID:    testutil.TestOrgID,
		DeviceID: 999,
	})

	if !errors.Is(err, providers.ErrNotFound) {
		t.Errorf("expected ErrNotFound, got: %v", err)
	}
}
