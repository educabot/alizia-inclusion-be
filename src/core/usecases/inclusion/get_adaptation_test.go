package inclusion_test

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"

	"github.com/educabot/alizia-inclusion-be/src/core/entities"
	"github.com/educabot/alizia-inclusion-be/src/core/providers"
	"github.com/educabot/alizia-inclusion-be/src/core/providers/mocks"
	"github.com/educabot/alizia-inclusion-be/src/core/usecases/inclusion"
	"github.com/educabot/alizia-inclusion-be/src/testutil"
)

func TestGetAdaptation(t *testing.T) {
	ctx := context.Background()

	t.Run("returns adaptation", func(t *testing.T) {
		// Arrange
		expected := testutil.NewAdaptation(1, 1, 1)
		mock := &mocks.MockAdaptationProvider{
			GetFn: func(_ context.Context, orgID uuid.UUID, id int64) (*entities.Adaptation, error) {
				return &expected, nil
			},
		}

		req := inclusion.GetAdaptationRequest{OrgID: testutil.TestOrgID, AdaptationID: 1}

		// Act
		got, err := inclusion.NewGetAdaptation(mock).Execute(ctx, req)

		// Assert
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if got == nil {
			t.Fatal("expected adaptation, got nil")
		}
		if got.ID != expected.ID {
			t.Errorf("got ID %d, want %d", got.ID, expected.ID)
		}
	})

	t.Run("rejects nil org_id", func(t *testing.T) {
		// Arrange
		called := false
		mock := &mocks.MockAdaptationProvider{
			GetFn: func(_ context.Context, _ uuid.UUID, _ int64) (*entities.Adaptation, error) {
				called = true
				return nil, nil
			},
		}

		req := inclusion.GetAdaptationRequest{OrgID: uuid.Nil, AdaptationID: 1}

		// Act
		_, err := inclusion.NewGetAdaptation(mock).Execute(ctx, req)

		// Assert
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

	t.Run("rejects zero adaptation_id", func(t *testing.T) {
		// Arrange
		called := false
		mock := &mocks.MockAdaptationProvider{
			GetFn: func(_ context.Context, _ uuid.UUID, _ int64) (*entities.Adaptation, error) {
				called = true
				return nil, nil
			},
		}

		req := inclusion.GetAdaptationRequest{OrgID: testutil.TestOrgID, AdaptationID: 0}

		// Act
		_, err := inclusion.NewGetAdaptation(mock).Execute(ctx, req)

		// Assert
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
		// Arrange
		mock := &mocks.MockAdaptationProvider{
			GetFn: func(_ context.Context, _ uuid.UUID, _ int64) (*entities.Adaptation, error) {
				return nil, errAdaptationNotFound
			},
		}

		req := inclusion.GetAdaptationRequest{OrgID: testutil.TestOrgID, AdaptationID: 99}

		// Act
		_, err := inclusion.NewGetAdaptation(mock).Execute(ctx, req)

		// Assert
		if err == nil {
			t.Error("expected error, got nil")
		}
		if !errors.Is(err, providers.ErrNotFound) {
			t.Errorf("expected ErrNotFound, got: %v", err)
		}
	})
}
