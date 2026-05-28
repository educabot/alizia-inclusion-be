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

func TestCreateAdaptation_CreatesAdaptationWithoutDevices(t *testing.T) {
	ctx := context.Background()

	mock := &mocks.MockAdaptationProvider{
		CreateFn: func(_ context.Context, a *entities.Adaptation) error {
			a.ID = 1
			return nil
		},
		GetFn: func(_ context.Context, orgID uuid.UUID, id int64) (*entities.Adaptation, error) {
			a := testutil.NewAdaptation(id, 1, 1)
			return &a, nil
		},
	}

	req := inclusion.CreateAdaptationRequest{
		OrgID:          testutil.TestOrgID,
		StudentID:      1,
		TeacherID:      1,
		Subject:        "Matematicas",
		AdaptationType: "actividad_adaptada",
	}

	got, err := inclusion.NewCreateAdaptation(mock).Execute(ctx, req)

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if got == nil {
		t.Fatal("expected adaptation, got nil")
	}
	if got.ID != 1 {
		t.Errorf("got ID %d, want 1", got.ID)
	}
}

func TestCreateAdaptation_CreatesAdaptationWithDevices(t *testing.T) {
	ctx := context.Background()

	setDevicesCalled := false
	mock := &mocks.MockAdaptationProvider{
		CreateFn: func(_ context.Context, a *entities.Adaptation) error {
			a.ID = 1
			return nil
		},
		SetDevicesFn: func(_ context.Context, adaptationID int64, deviceIDs []int64) error {
			setDevicesCalled = true
			return nil
		},
		GetFn: func(_ context.Context, orgID uuid.UUID, id int64) (*entities.Adaptation, error) {
			a := testutil.NewAdaptation(id, 1, 1)
			return &a, nil
		},
	}

	req := inclusion.CreateAdaptationRequest{
		OrgID:          testutil.TestOrgID,
		StudentID:      1,
		TeacherID:      1,
		Subject:        "Matematicas",
		AdaptationType: "actividad_adaptada",
		DeviceIDs:      []int64{10, 20},
	}

	got, err := inclusion.NewCreateAdaptation(mock).Execute(ctx, req)

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if got == nil {
		t.Fatal("expected adaptation, got nil")
	}
	if !setDevicesCalled {
		t.Error("expected SetDevices to be called, but it was not")
	}
}

func TestCreateAdaptation_DefaultsAdaptationTypeAndPersistsTitleWhenTypeOmitted(t *testing.T) {
	ctx := context.Background()

	var captured *entities.Adaptation
	mock := &mocks.MockAdaptationProvider{
		CreateFn: func(_ context.Context, a *entities.Adaptation) error {
			captured = a
			a.ID = 1
			return nil
		},
		GetFn: func(_ context.Context, _ uuid.UUID, id int64) (*entities.Adaptation, error) {
			a := testutil.NewAdaptation(id, 1, 1)
			return &a, nil
		},
	}

	req := inclusion.CreateAdaptationRequest{
		OrgID:          testutil.TestOrgID,
		StudentID:      1,
		TeacherID:      1,
		Title:          "Secuencia con apoyos visuales",
		Subject:        "Matematicas",
		AdaptationType: "",
	}

	_, err := inclusion.NewCreateAdaptation(mock).Execute(ctx, req)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if captured == nil {
		t.Fatal("expected Create to be called")
	}
	if captured.AdaptationType != "actividad_adaptada" {
		t.Errorf("got AdaptationType %q, want %q", captured.AdaptationType, "actividad_adaptada")
	}
	if captured.Title != "Secuencia con apoyos visuales" {
		t.Errorf("got Title %q, want %q", captured.Title, "Secuencia con apoyos visuales")
	}
}

func TestCreateAdaptation_RejectsNilOrgID(t *testing.T) {
	ctx := context.Background()

	called := false
	mock := &mocks.MockAdaptationProvider{
		CreateFn: func(_ context.Context, _ *entities.Adaptation) error {
			called = true
			return nil
		},
	}

	req := inclusion.CreateAdaptationRequest{
		OrgID:     uuid.Nil,
		StudentID: 1,
		TeacherID: 1,
		Subject:   "Matematicas",
	}

	_, err := inclusion.NewCreateAdaptation(mock).Execute(ctx, req)

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

func TestCreateAdaptation_RejectsZeroStudentID(t *testing.T) {
	ctx := context.Background()

	called := false
	mock := &mocks.MockAdaptationProvider{
		CreateFn: func(_ context.Context, _ *entities.Adaptation) error {
			called = true
			return nil
		},
	}

	req := inclusion.CreateAdaptationRequest{
		OrgID:     testutil.TestOrgID,
		StudentID: 0,
		TeacherID: 1,
		Subject:   "Matematicas",
	}

	_, err := inclusion.NewCreateAdaptation(mock).Execute(ctx, req)

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

func TestCreateAdaptation_RejectsZeroTeacherID(t *testing.T) {
	ctx := context.Background()

	called := false
	mock := &mocks.MockAdaptationProvider{
		CreateFn: func(_ context.Context, _ *entities.Adaptation) error {
			called = true
			return nil
		},
	}

	req := inclusion.CreateAdaptationRequest{
		OrgID:     testutil.TestOrgID,
		StudentID: 1,
		TeacherID: 0,
		Subject:   "Matematicas",
	}

	_, err := inclusion.NewCreateAdaptation(mock).Execute(ctx, req)

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

func TestCreateAdaptation_RejectsEmptySubject(t *testing.T) {
	ctx := context.Background()

	called := false
	mock := &mocks.MockAdaptationProvider{
		CreateFn: func(_ context.Context, _ *entities.Adaptation) error {
			called = true
			return nil
		},
	}

	req := inclusion.CreateAdaptationRequest{
		OrgID:     testutil.TestOrgID,
		StudentID: 1,
		TeacherID: 1,
		Subject:   "",
	}

	_, err := inclusion.NewCreateAdaptation(mock).Execute(ctx, req)

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
