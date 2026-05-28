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

func TestGetRamp_ReturnsRamp(t *testing.T) {
	ctx := context.Background()
	expected := testutil.NewRamp(1, "Ramp 1")
	mock := &mocks.MockRampProvider{
		GetRampFn: func(_ context.Context, orgID uuid.UUID, id int64) (*entities.Ramp, error) {
			return &expected, nil
		},
	}

	got, err := catalog.NewGetRamp(mock).Execute(ctx, catalog.GetRampRequest{OrgID: testutil.TestOrgID, RampID: 1})

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if got == nil {
		t.Fatal("expected ramp, got nil")
	}
	if got.ID != expected.ID || got.Name != expected.Name {
		t.Errorf("got {ID:%d Name:%q}, want {ID:%d Name:%q}",
			got.ID, got.Name, expected.ID, expected.Name)
	}
}

func TestGetRamp_RejectsNilOrgID(t *testing.T) {
	ctx := context.Background()
	called := false
	mock := &mocks.MockRampProvider{
		GetRampFn: func(_ context.Context, _ uuid.UUID, _ int64) (*entities.Ramp, error) {
			called = true
			return nil, nil
		},
	}

	_, err := catalog.NewGetRamp(mock).Execute(ctx, catalog.GetRampRequest{OrgID: uuid.Nil, RampID: 1})

	if err == nil {
		t.Error("expected validation error, got nil")
	}
	if !errors.Is(err, providers.ErrValidation) {
		t.Errorf("expected ErrValidation, got: %v", err)
	}
	if called {
		t.Error("mock should not have been called for invalid request")
	}
}

func TestGetRamp_RejectsZeroRampID(t *testing.T) {
	ctx := context.Background()
	called := false
	mock := &mocks.MockRampProvider{
		GetRampFn: func(_ context.Context, _ uuid.UUID, _ int64) (*entities.Ramp, error) {
			called = true
			return nil, nil
		},
	}

	_, err := catalog.NewGetRamp(mock).Execute(ctx, catalog.GetRampRequest{OrgID: testutil.TestOrgID, RampID: 0})

	if err == nil {
		t.Error("expected validation error, got nil")
	}
	if !errors.Is(err, providers.ErrValidation) {
		t.Errorf("expected ErrValidation, got: %v", err)
	}
	if called {
		t.Error("mock should not have been called for invalid request")
	}
}

func TestGetRamp_ReturnsNotFound(t *testing.T) {
	ctx := context.Background()
	mock := &mocks.MockRampProvider{
		GetRampFn: func(_ context.Context, _ uuid.UUID, _ int64) (*entities.Ramp, error) {
			return nil, errRampNotFound
		},
	}

	_, err := catalog.NewGetRamp(mock).Execute(ctx, catalog.GetRampRequest{OrgID: testutil.TestOrgID, RampID: 99})

	if err == nil {
		t.Error("expected error, got nil")
	}
	if !errors.Is(err, providers.ErrNotFound) {
		t.Errorf("expected ErrNotFound, got: %v", err)
	}
}
