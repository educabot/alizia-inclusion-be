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

func TestGetAdaptation_ReturnsAdaptation(t *testing.T) {
	ctx := context.Background()
	expected := testutil.NewAdaptation(1, 1, 1)
	mock := &mocks.MockAdaptationProvider{
		GetFn: func(_ context.Context, orgID uuid.UUID, id int64) (*entities.Adaptation, error) {
			return &expected, nil
		},
	}

	req := inclusion.GetAdaptationRequest{OrgID: testutil.TestOrgID, AdaptationID: 1}

	got, err := inclusion.NewGetAdaptation(mock).Execute(ctx, req)

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if got == nil {
		t.Fatal("expected adaptation, got nil")
	}
	if got.ID != expected.ID {
		t.Errorf("got ID %d, want %d", got.ID, expected.ID)
	}
}

func TestGetAdaptation_RejectsNilOrgID(t *testing.T) {
	ctx := context.Background()
	called := false
	mock := &mocks.MockAdaptationProvider{
		GetFn: func(_ context.Context, _ uuid.UUID, _ int64) (*entities.Adaptation, error) {
			called = true
			return nil, nil
		},
	}

	req := inclusion.GetAdaptationRequest{OrgID: uuid.Nil, AdaptationID: 1}

	_, err := inclusion.NewGetAdaptation(mock).Execute(ctx, req)

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

func TestGetAdaptation_RejectsZeroAdaptationID(t *testing.T) {
	ctx := context.Background()
	called := false
	mock := &mocks.MockAdaptationProvider{
		GetFn: func(_ context.Context, _ uuid.UUID, _ int64) (*entities.Adaptation, error) {
			called = true
			return nil, nil
		},
	}

	req := inclusion.GetAdaptationRequest{OrgID: testutil.TestOrgID, AdaptationID: 0}

	_, err := inclusion.NewGetAdaptation(mock).Execute(ctx, req)

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

func TestGetAdaptation_ReturnsNotFound(t *testing.T) {
	ctx := context.Background()
	mock := &mocks.MockAdaptationProvider{
		GetFn: func(_ context.Context, _ uuid.UUID, _ int64) (*entities.Adaptation, error) {
			return nil, errAdaptationNotFound
		},
	}

	req := inclusion.GetAdaptationRequest{OrgID: testutil.TestOrgID, AdaptationID: 99}

	_, err := inclusion.NewGetAdaptation(mock).Execute(ctx, req)

	if err == nil {
		t.Error("expected error, got nil")
	}
	if !errors.Is(err, providers.ErrNotFound) {
		t.Errorf("expected ErrNotFound, got: %v", err)
	}
}
