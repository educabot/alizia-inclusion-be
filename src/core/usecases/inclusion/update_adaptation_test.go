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

func TestUpdateAdaptation_UpdatesAdaptation(t *testing.T) {
	ctx := context.Background()
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

	got, err := inclusion.NewUpdateAdaptation(mock).Execute(ctx, req)

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
}

func TestUpdateAdaptation_AppliesTitleWhenProvided(t *testing.T) {
	ctx := context.Background()
	existing := testutil.NewAdaptation(1, 1, 1)
	var captured *entities.Adaptation
	mock := &mocks.MockAdaptationProvider{
		GetFn: func(_ context.Context, _ uuid.UUID, _ int64) (*entities.Adaptation, error) {
			return &existing, nil
		},
		UpdateFn: func(_ context.Context, a *entities.Adaptation) error {
			captured = a
			return nil
		},
	}

	req := inclusion.UpdateAdaptationRequest{
		OrgID:        testutil.TestOrgID,
		AdaptationID: 1,
		Title:        testutil.Ptr("Secuencia con apoyos visuales"),
	}

	_, err := inclusion.NewUpdateAdaptation(mock).Execute(ctx, req)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if captured == nil {
		t.Fatal("expected Update to be called")
	}
	if captured.Title != "Secuencia con apoyos visuales" {
		t.Errorf("got Title %q, want %q", captured.Title, "Secuencia con apoyos visuales")
	}
}

func TestUpdateAdaptation_RejectsNilOrgID(t *testing.T) {
	ctx := context.Background()
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

	_, err := inclusion.NewUpdateAdaptation(mock).Execute(ctx, req)

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

func TestUpdateAdaptation_RejectsZeroAdaptationID(t *testing.T) {
	ctx := context.Background()
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

	_, err := inclusion.NewUpdateAdaptation(mock).Execute(ctx, req)

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

func TestUpdateAdaptation_RejectsInvalidStatus(t *testing.T) {
	ctx := context.Background()
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

	_, err := inclusion.NewUpdateAdaptation(mock).Execute(ctx, req)

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

func TestUpdateAdaptation_ReturnsNotFound(t *testing.T) {
	ctx := context.Background()
	mock := &mocks.MockAdaptationProvider{
		GetFn: func(_ context.Context, _ uuid.UUID, _ int64) (*entities.Adaptation, error) {
			return nil, errAdaptationNotFound
		},
	}

	req := inclusion.UpdateAdaptationRequest{
		OrgID:        testutil.TestOrgID,
		AdaptationID: 99,
	}

	_, err := inclusion.NewUpdateAdaptation(mock).Execute(ctx, req)

	if err == nil {
		t.Error("expected error, got nil")
	}
	if !errors.Is(err, providers.ErrNotFound) {
		t.Errorf("expected ErrNotFound, got: %v", err)
	}
}
