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

func TestUpdateAdaptation(t *testing.T) {
	ctx := context.Background()

	t.Run("updates adaptation", func(t *testing.T) {
		// Arrange
		getCallCount := 0
		existing := testutil.NewAdaptation(1, 1, 1)
		updated := testutil.NewAdaptation(1, 1, 1)
		updated.Subject = "Lengua"
		updated.Status = "probado"

		mock := &mocks.MockAdaptationProvider{
			GetFn: func(_ context.Context, orgID uuid.UUID, id int64) (*entities.Adaptation, error) {
				getCallCount++
				if getCallCount == 1 {
					return &existing, nil
				}
				return &updated, nil
			},
			UpdateFn: func(_ context.Context, a *entities.Adaptation) error {
				return nil
			},
		}

		req := inclusion.UpdateAdaptationRequest{
			OrgID:        testutil.TestOrgID,
			AdaptationID: 1,
			Subject:      testutil.Ptr("Lengua"),
			Status:       testutil.Ptr("probado"),
		}

		// Act
		got, err := inclusion.NewUpdateAdaptation(mock).Execute(ctx, req)

		// Assert
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if got == nil {
			t.Fatal("expected adaptation, got nil")
		}
		if got.Subject != "Lengua" {
			t.Errorf("got Subject %q, want %q", got.Subject, "Lengua")
		}
		if got.Status != "probado" {
			t.Errorf("got Status %q, want %q", got.Status, "probado")
		}
		if getCallCount != 2 {
			t.Errorf("expected Get to be called 2 times, got %d", getCallCount)
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

		req := inclusion.UpdateAdaptationRequest{
			OrgID:        uuid.Nil,
			AdaptationID: 1,
		}

		// Act
		_, err := inclusion.NewUpdateAdaptation(mock).Execute(ctx, req)

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

		req := inclusion.UpdateAdaptationRequest{
			OrgID:        testutil.TestOrgID,
			AdaptationID: 0,
		}

		// Act
		_, err := inclusion.NewUpdateAdaptation(mock).Execute(ctx, req)

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

	t.Run("rejects invalid status", func(t *testing.T) {
		// Arrange
		called := false
		mock := &mocks.MockAdaptationProvider{
			GetFn: func(_ context.Context, _ uuid.UUID, _ int64) (*entities.Adaptation, error) {
				called = true
				return nil, nil
			},
		}

		req := inclusion.UpdateAdaptationRequest{
			OrgID:        testutil.TestOrgID,
			AdaptationID: 1,
			Status:       testutil.Ptr("invalid"),
		}

		// Act
		_, err := inclusion.NewUpdateAdaptation(mock).Execute(ctx, req)

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

		req := inclusion.UpdateAdaptationRequest{
			OrgID:        testutil.TestOrgID,
			AdaptationID: 99,
		}

		// Act
		_, err := inclusion.NewUpdateAdaptation(mock).Execute(ctx, req)

		// Assert
		if err == nil {
			t.Error("expected error, got nil")
		}
		if !errors.Is(err, providers.ErrNotFound) {
			t.Errorf("expected ErrNotFound, got: %v", err)
		}
	})
}
