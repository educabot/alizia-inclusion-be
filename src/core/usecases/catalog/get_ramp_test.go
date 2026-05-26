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

func TestGetRamp(t *testing.T) {
	ctx := context.Background()

	t.Run("returns ramp", func(t *testing.T) {
		expected := testutil.NewRamp(1, "Ramp 1")

		mock := &mocks.MockRampProvider{
			GetRampFn: func(_ context.Context, orgID uuid.UUID, id int64) (*entities.Ramp, error) {
				return &expected, nil
			},
		}

		req := catalog.GetRampRequest{OrgID: testutil.TestOrgID, RampID: 1}
		got, err := catalog.NewGetRamp(mock).Execute(ctx, req)
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
	})

	t.Run("rejects nil org_id", func(t *testing.T) {
		called := false
		mock := &mocks.MockRampProvider{
			GetRampFn: func(_ context.Context, _ uuid.UUID, _ int64) (*entities.Ramp, error) {
				called = true
				return nil, nil
			},
		}

		req := catalog.GetRampRequest{OrgID: uuid.Nil, RampID: 1}
		_, err := catalog.NewGetRamp(mock).Execute(ctx, req)
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

	t.Run("rejects zero ramp_id", func(t *testing.T) {
		called := false
		mock := &mocks.MockRampProvider{
			GetRampFn: func(_ context.Context, _ uuid.UUID, _ int64) (*entities.Ramp, error) {
				called = true
				return nil, nil
			},
		}

		req := catalog.GetRampRequest{OrgID: testutil.TestOrgID, RampID: 0}
		_, err := catalog.NewGetRamp(mock).Execute(ctx, req)
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

	t.Run("returns not found", func(t *testing.T) {
		mock := &mocks.MockRampProvider{
			GetRampFn: func(_ context.Context, _ uuid.UUID, _ int64) (*entities.Ramp, error) {
				return nil, errRampNotFound
			},
		}

		req := catalog.GetRampRequest{OrgID: testutil.TestOrgID, RampID: 99}
		_, err := catalog.NewGetRamp(mock).Execute(ctx, req)
		if err == nil {
			t.Error("expected error, got nil")
		}
		if !errors.Is(err, providers.ErrNotFound) {
			t.Errorf("expected ErrNotFound, got: %v", err)
		}
	})
}
