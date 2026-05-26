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

func TestListTeachers(t *testing.T) {
	ctx := context.Background()

	t.Run("returns teachers", func(t *testing.T) {
		teachers := []entities.User{
			testutil.NewUser(1, "Ana Garcia"),
			testutil.NewUser(2, "Luis Perez"),
		}
		var calledRole string
		mock := &mocks.MockUserProvider{
			ListByRoleFn: func(ctx context.Context, orgID uuid.UUID, role string) ([]entities.User, error) {
				calledRole = role
				return teachers, nil
			},
		}

		uc := management.NewListTeachers(mock)
		got, err := uc.Execute(ctx, management.ListTeachersRequest{OrgID: testutil.TestOrgID})
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if calledRole != "teacher" {
			t.Errorf("expected role %q, got %q", "teacher", calledRole)
		}
		if len(got) != len(teachers) {
			t.Fatalf("expected %d teachers, got %d", len(teachers), len(got))
		}
		if got[0].ID != teachers[0].ID {
			t.Errorf("expected first teacher ID %d, got %d", teachers[0].ID, got[0].ID)
		}
	})

	t.Run("rejects nil org_id", func(t *testing.T) {
		mock := &mocks.MockUserProvider{}

		uc := management.NewListTeachers(mock)
		_, err := uc.Execute(ctx, management.ListTeachersRequest{OrgID: uuid.Nil})
		if err == nil {
			t.Fatal("expected validation error, got nil")
		}
		if !errors.Is(err, providers.ErrValidation) {
			t.Errorf("expected ErrValidation, got %v", err)
		}
	})
}
