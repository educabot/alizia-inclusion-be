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

func TestListAdaptations_ReturnsAllAdaptations(t *testing.T) {
	ctx := context.Background()
	expected := []entities.Adaptation{
		testutil.NewAdaptation(1, 1, 1),
		testutil.NewAdaptation(2, 2, 1),
	}
	mock := &mocks.MockAdaptationProvider{
		ListFn: func(_ context.Context, orgID uuid.UUID, studentID *int64) ([]entities.Adaptation, error) {
			return expected, nil
		},
	}

	req := inclusion.ListAdaptationsRequest{OrgID: testutil.TestOrgID, StudentID: nil}

	got, err := inclusion.NewListAdaptations(mock).Execute(ctx, req)

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if len(got) != len(expected) {
		t.Errorf("got %d adaptations, want %d", len(got), len(expected))
	}
	for i, a := range got {
		if a.ID != expected[i].ID {
			t.Errorf("adaptation[%d].ID = %d, want %d", i, a.ID, expected[i].ID)
		}
	}
}

func TestListAdaptations_FiltersByStudent(t *testing.T) {
	ctx := context.Background()
	wantStudentID := int64(1)
	var capturedStudentID *int64
	expected := []entities.Adaptation{
		testutil.NewAdaptation(1, wantStudentID, 1),
	}
	mock := &mocks.MockAdaptationProvider{
		ListFn: func(_ context.Context, _ uuid.UUID, studentID *int64) ([]entities.Adaptation, error) {
			capturedStudentID = studentID
			return expected, nil
		},
	}

	req := inclusion.ListAdaptationsRequest{
		OrgID:     testutil.TestOrgID,
		StudentID: testutil.Ptr(wantStudentID),
	}

	got, err := inclusion.NewListAdaptations(mock).Execute(ctx, req)

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if len(got) != 1 {
		t.Errorf("got %d adaptations, want 1", len(got))
	}
	if capturedStudentID == nil {
		t.Fatal("studentID was not passed to mock")
	}
	if *capturedStudentID != wantStudentID {
		t.Errorf("mock received studentID %d, want %d", *capturedStudentID, wantStudentID)
	}
}

func TestListAdaptations_RejectsNilOrgID(t *testing.T) {
	ctx := context.Background()
	called := false
	mock := &mocks.MockAdaptationProvider{
		ListFn: func(_ context.Context, _ uuid.UUID, _ *int64) ([]entities.Adaptation, error) {
			called = true
			return nil, nil
		},
	}

	req := inclusion.ListAdaptationsRequest{OrgID: uuid.Nil}

	_, err := inclusion.NewListAdaptations(mock).Execute(ctx, req)

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
