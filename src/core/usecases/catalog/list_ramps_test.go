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

func TestListRamps_ReturnsRampsForValidOrg(t *testing.T) {
	ctx := context.Background()
	expected := []entities.Ramp{
		testutil.NewRamp(1, "Ramp 1"),
		testutil.NewRamp(2, "Ramp 2"),
	}
	mock := &mocks.MockRampProvider{
		ListRampsFn: func(_ context.Context, orgID uuid.UUID) ([]entities.Ramp, error) {
			return expected, nil
		},
	}

	got, err := catalog.NewListRamps(mock).Execute(ctx, catalog.ListRampsRequest{OrgID: testutil.TestOrgID})

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if len(got) != len(expected) {
		t.Errorf("got %d ramps, want %d", len(got), len(expected))
	}
	for i, r := range got {
		if r.ID != expected[i].ID || r.Name != expected[i].Name {
			t.Errorf("ramp[%d] = {ID:%d Name:%q}, want {ID:%d Name:%q}",
				i, r.ID, r.Name, expected[i].ID, expected[i].Name)
		}
	}
}

func TestListRamps_RejectsNilOrgID(t *testing.T) {
	ctx := context.Background()
	called := false
	mock := &mocks.MockRampProvider{
		ListRampsFn: func(_ context.Context, _ uuid.UUID) ([]entities.Ramp, error) {
			called = true
			return nil, nil
		},
	}

	_, err := catalog.NewListRamps(mock).Execute(ctx, catalog.ListRampsRequest{OrgID: uuid.Nil})

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

func TestListRamps_PropagatesProviderError(t *testing.T) {
	ctx := context.Background()
	mock := &mocks.MockRampProvider{
		ListRampsFn: func(_ context.Context, _ uuid.UUID) ([]entities.Ramp, error) {
			return nil, errDB
		},
	}

	_, err := catalog.NewListRamps(mock).Execute(ctx, catalog.ListRampsRequest{OrgID: testutil.TestOrgID})

	if err == nil {
		t.Error("expected error, got nil")
	}
	if !errors.Is(err, errDB) {
		t.Errorf("expected errDB, got: %v", err)
	}
}
