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

func TestCreateAdaptation(t *testing.T) {
	ctx := context.Background()

	t.Run("creates adaptation without devices", func(t *testing.T) {
		// Arrange
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

		// Act
		got, err := inclusion.NewCreateAdaptation(mock).Execute(ctx, req)

		// Assert
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if got == nil {
			t.Fatal("expected adaptation, got nil")
		}
		if got.ID != 1 {
			t.Errorf("got ID %d, want 1", got.ID)
		}
	})

	t.Run("creates adaptation with devices", func(t *testing.T) {
		// Arrange
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

		// Act
		got, err := inclusion.NewCreateAdaptation(mock).Execute(ctx, req)

		// Assert
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if got == nil {
			t.Fatal("expected adaptation, got nil")
		}
		if !setDevicesCalled {
			t.Error("expected SetDevices to be called, but it was not")
		}
	})

	t.Run("rejects nil org_id", func(t *testing.T) {
		// Arrange
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

		// Act
		_, err := inclusion.NewCreateAdaptation(mock).Execute(ctx, req)

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

	t.Run("rejects zero student_id", func(t *testing.T) {
		// Arrange
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

		// Act
		_, err := inclusion.NewCreateAdaptation(mock).Execute(ctx, req)

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

	t.Run("rejects zero teacher_id", func(t *testing.T) {
		// Arrange
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

		// Act
		_, err := inclusion.NewCreateAdaptation(mock).Execute(ctx, req)

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

	t.Run("rejects empty subject", func(t *testing.T) {
		// Arrange
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

		// Act
		_, err := inclusion.NewCreateAdaptation(mock).Execute(ctx, req)

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
}
