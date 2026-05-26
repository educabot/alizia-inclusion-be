package auth_test

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"

	"github.com/educabot/alizia-inclusion-be/src/core/entities"
	"github.com/educabot/alizia-inclusion-be/src/core/providers"
	"github.com/educabot/alizia-inclusion-be/src/core/providers/mocks"
	"github.com/educabot/alizia-inclusion-be/src/core/usecases/auth"
	"github.com/educabot/alizia-inclusion-be/src/testutil"
)

func TestGetMe(t *testing.T) {
	ctx := context.Background()

	t.Run("returns user", func(t *testing.T) {
		// Arrange
		want := testutil.NewUser(42, "Ana")
		mock := &mocks.MockUserProvider{
			GetByIDFn: func(_ context.Context, _ uuid.UUID, _ int64) (*entities.User, error) {
				return &want, nil
			},
		}
		uc := auth.NewGetMe(mock)

		// Act
		got, err := uc.Execute(ctx, auth.GetMeRequest{
			OrgID:  testutil.TestOrgID,
			UserID: 42,
		})

		// Assert
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if got == nil {
			t.Fatal("expected non-nil user, got nil")
		}
		if got.ID != want.ID {
			t.Errorf("user ID = %d, want %d", got.ID, want.ID)
		}
		if got.Name != want.Name {
			t.Errorf("user Name = %q, want %q", got.Name, want.Name)
		}
	})

	t.Run("rejects nil org_id", func(t *testing.T) {
		// Arrange
		mock := &mocks.MockUserProvider{}
		uc := auth.NewGetMe(mock)

		// Act
		_, err := uc.Execute(ctx, auth.GetMeRequest{
			OrgID:  uuid.Nil,
			UserID: 1,
		})

		// Assert
		if err == nil {
			t.Fatal("expected validation error, got nil")
		}
		if !errors.Is(err, providers.ErrValidation) {
			t.Errorf("expected ErrValidation, got %v", err)
		}
	})

	t.Run("rejects zero user_id", func(t *testing.T) {
		// Arrange
		mock := &mocks.MockUserProvider{}
		uc := auth.NewGetMe(mock)

		// Act
		_, err := uc.Execute(ctx, auth.GetMeRequest{
			OrgID:  testutil.TestOrgID,
			UserID: 0,
		})

		// Assert
		if err == nil {
			t.Fatal("expected validation error, got nil")
		}
		if !errors.Is(err, providers.ErrValidation) {
			t.Errorf("expected ErrValidation, got %v", err)
		}
	})

	t.Run("returns not found", func(t *testing.T) {
		// Arrange
		mock := &mocks.MockUserProvider{
			GetByIDFn: func(_ context.Context, _ uuid.UUID, _ int64) (*entities.User, error) {
				return nil, errUserNotFound
			},
		}
		uc := auth.NewGetMe(mock)

		// Act
		got, err := uc.Execute(ctx, auth.GetMeRequest{
			OrgID:  testutil.TestOrgID,
			UserID: 99,
		})

		// Assert
		if !errors.Is(err, providers.ErrNotFound) {
			t.Errorf("expected ErrNotFound, got %v", err)
		}
		if got != nil {
			t.Errorf("expected nil user, got %+v", got)
		}
	})
}
