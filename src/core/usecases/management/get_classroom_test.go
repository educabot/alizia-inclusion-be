package management_test

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"

	"github.com/educabot/alizia-inclusion-be/src/core/entities"
	"github.com/educabot/alizia-inclusion-be/src/core/providers"
	"github.com/educabot/alizia-inclusion-be/src/core/providers/mocks"
	"github.com/educabot/alizia-inclusion-be/src/core/usecases/management"
	"github.com/educabot/alizia-inclusion-be/src/testutil"
)

func TestGetClassroom(t *testing.T) {
	ctx := context.Background()

	t.Run("returns classroom", func(t *testing.T) {
		classroom := testutil.NewClassroom(1, "1A")
		mock := &mocks.MockClassroomProvider{
			GetFn: func(ctx context.Context, orgID uuid.UUID, id int64) (*entities.Classroom, error) {
				return &classroom, nil
			},
		}

		uc := management.NewGetClassroom(mock)
		got, err := uc.Execute(ctx, management.GetClassroomRequest{
			OrgID:       testutil.TestOrgID,
			ClassroomID: 1,
		})
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if got.ID != classroom.ID {
			t.Errorf("expected classroom ID %d, got %d", classroom.ID, got.ID)
		}
		if got.Name != classroom.Name {
			t.Errorf("expected classroom name %q, got %q", classroom.Name, got.Name)
		}
	})

	t.Run("rejects nil org_id", func(t *testing.T) {
		mock := &mocks.MockClassroomProvider{}

		uc := management.NewGetClassroom(mock)
		_, err := uc.Execute(ctx, management.GetClassroomRequest{
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

		uc := management.NewGetClassroom(mock)
		_, err := uc.Execute(ctx, management.GetClassroomRequest{
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
			GetFn: func(ctx context.Context, orgID uuid.UUID, id int64) (*entities.Classroom, error) {
				return nil, errClassroomNotFound
			},
		}

		uc := management.NewGetClassroom(mock)
		_, err := uc.Execute(ctx, management.GetClassroomRequest{
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
