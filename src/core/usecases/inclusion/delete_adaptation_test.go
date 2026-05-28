package inclusion_test

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"

	"github.com/educabot/alizia-inclusion-be/src/core/providers"
	"github.com/educabot/alizia-inclusion-be/src/core/providers/mocks"
	"github.com/educabot/alizia-inclusion-be/src/core/usecases/inclusion"
	"github.com/educabot/alizia-inclusion-be/src/testutil"
)

func TestDeleteAdaptation_DeletesAdaptation(t *testing.T) {
	ctx := context.Background()
	deleteCalled := false
	mock := &mocks.MockAdaptationProvider{
		DeleteFn: func(_ context.Context, orgID uuid.UUID, id int64) error {
			deleteCalled = true
			return nil
		},
	}

	req := inclusion.DeleteAdaptationRequest{
		OrgID:        testutil.TestOrgID,
		AdaptationID: 1,
	}

	err := inclusion.NewDeleteAdaptation(mock).Execute(ctx, req)

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if !deleteCalled {
		t.Error("expected Delete to be called, but it was not")
	}
}

func TestDeleteAdaptation_RejectsNilOrgID(t *testing.T) {
	ctx := context.Background()
	called := false
	mock := &mocks.MockAdaptationProvider{
		DeleteFn: func(_ context.Context, _ uuid.UUID, _ int64) error {
			called = true
			return nil
		},
	}

	req := inclusion.DeleteAdaptationRequest{
		OrgID:        uuid.Nil,
		AdaptationID: 1,
	}

	err := inclusion.NewDeleteAdaptation(mock).Execute(ctx, req)

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

func TestDeleteAdaptation_RejectsZeroAdaptationID(t *testing.T) {
	ctx := context.Background()
	called := false
	mock := &mocks.MockAdaptationProvider{
		DeleteFn: func(_ context.Context, _ uuid.UUID, _ int64) error {
			called = true
			return nil
		},
	}

	req := inclusion.DeleteAdaptationRequest{
		OrgID:        testutil.TestOrgID,
		AdaptationID: 0,
	}

	err := inclusion.NewDeleteAdaptation(mock).Execute(ctx, req)

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
