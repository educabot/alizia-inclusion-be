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

func TestGetMe_ReturnsUser(t *testing.T) {
	ctx := context.Background()
	want := testutil.NewUser(42, "Ana")
	mock := &mocks.MockUserProvider{
		GetByIDFn: func(_ context.Context, _ uuid.UUID, _ int64) (*entities.User, error) {
			return &want, nil
		},
	}
	uc := auth.NewGetMe(mock)

	got, err := uc.Execute(ctx, auth.GetMeRequest{
		OrgID:  testutil.TestOrgID,
		UserID: 42,
	})

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
}

func TestGetMe_RejectsNilOrgID(t *testing.T) {
	ctx := context.Background()
	mock := &mocks.MockUserProvider{}
	uc := auth.NewGetMe(mock)

	_, err := uc.Execute(ctx, auth.GetMeRequest{
		OrgID:  uuid.Nil,
		UserID: 1,
	})

	if err == nil {
		t.Fatal("expected validation error, got nil")
	}
	if !errors.Is(err, providers.ErrValidation) {
		t.Errorf("expected ErrValidation, got %v", err)
	}
}

func TestGetMe_RejectsZeroUserID(t *testing.T) {
	ctx := context.Background()
	mock := &mocks.MockUserProvider{}
	uc := auth.NewGetMe(mock)

	_, err := uc.Execute(ctx, auth.GetMeRequest{
		OrgID:  testutil.TestOrgID,
		UserID: 0,
	})

	if err == nil {
		t.Fatal("expected validation error, got nil")
	}
	if !errors.Is(err, providers.ErrValidation) {
		t.Errorf("expected ErrValidation, got %v", err)
	}
}

func TestGetMe_ReturnsNotFound(t *testing.T) {
	ctx := context.Background()
	mock := &mocks.MockUserProvider{
		GetByIDFn: func(_ context.Context, _ uuid.UUID, _ int64) (*entities.User, error) {
			return nil, errUserNotFound
		},
	}
	uc := auth.NewGetMe(mock)

	got, err := uc.Execute(ctx, auth.GetMeRequest{
		OrgID:  testutil.TestOrgID,
		UserID: 99,
	})

	if !errors.Is(err, providers.ErrNotFound) {
		t.Errorf("expected ErrNotFound, got %v", err)
	}
	if got != nil {
		t.Errorf("expected nil user, got %+v", got)
	}
}
