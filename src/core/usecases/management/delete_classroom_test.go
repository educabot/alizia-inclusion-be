package management_test

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"

	"github.com/educabot/alizia-inclusion-be/src/core/providers"
	"github.com/educabot/alizia-inclusion-be/src/core/providers/mocks"
	"github.com/educabot/alizia-inclusion-be/src/core/usecases/management"
	"github.com/educabot/alizia-inclusion-be/src/testutil"
)

func TestDeleteClassroom(t *testing.T) {
	ctx := context.Background()

	t.Run("deletes classroom", func(t *testing.T) {
		called := false
		mock := &mocks.MockClassroomProvider{
			DeleteFn: func(ctx context.Context, orgID uuid.UUID, id int64) error {
				called = true
				return nil
			},
		}

		uc := management.NewDeleteClassroom(mock)
		err := uc.Execute(ctx, management.DeleteClassroomRequest{
			OrgID:       testutil.TestOrgID,
			ClassroomID: 1,
		})
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if !called {
			t.Error("expected DeleteFn to be called, but it was not")
		}
	})

	t.Run("rejects nil org_id", func(t *testing.T) {
		mock := &mocks.MockClassroomProvider{}

		uc := management.NewDeleteClassroom(mock)
		err := uc.Execute(ctx, management.DeleteClassroomRequest{
			OrgID:       uuid.Nil,
			ClassroomID: 1,
		})
		if err == nil {
			t.Fatal("expected validation error, got nil")
		}
		if !errors.Is(err, providers.ErrValidation) {
			t.Errorf("expected ErrValidation, got %v", err)
		}
	})

	t.Run("rejects zero classroom_id", func(t *testing.T) {
		mock := &mocks.MockClassroomProvider{}

		uc := management.NewDeleteClassroom(mock)
		err := uc.Execute(ctx, management.DeleteClassroomRequest{
			OrgID:       testutil.TestOrgID,
			ClassroomID: 0,
		})
		if err == nil {
			t.Fatal("expected validation error, got nil")
		}
		if !errors.Is(err, providers.ErrValidation) {
			t.Errorf("expected ErrValidation, got %v", err)
		}
	})

	t.Run("returns not found", func(t *testing.T) {
		mock := &mocks.MockClassroomProvider{
			DeleteFn: func(ctx context.Context, orgID uuid.UUID, id int64) error {
				return errClassroomNotFound
			},
		}

		uc := management.NewDeleteClassroom(mock)
		err := uc.Execute(ctx, management.DeleteClassroomRequest{
			OrgID:       testutil.TestOrgID,
			ClassroomID: 99,
		})
		if err == nil {
			t.Fatal("expected not found error, got nil")
		}
		if !errors.Is(err, providers.ErrNotFound) {
			t.Errorf("expected ErrNotFound, got %v", err)
		}
	})
}
