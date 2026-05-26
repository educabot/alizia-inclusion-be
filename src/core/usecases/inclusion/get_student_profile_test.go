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

func TestGetStudentProfile(t *testing.T) {
	ctx := context.Background()

	t.Run("returns student with profile", func(t *testing.T) {
		expected := testutil.NewStudentWithProfile(1, 1, "Lucas", []string{"distraccion"})
		mock := &mocks.MockStudentProvider{
			GetStudentFn: func(_ context.Context, _ uuid.UUID, id int64) (*entities.Student, error) {
				s := expected
				return &s, nil
			},
		}

		got, err := inclusion.NewGetStudentProfile(mock).Execute(ctx, inclusion.GetStudentProfileRequest{
			OrgID:     testutil.TestOrgID,
			StudentID: 1,
		})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if got.Profile == nil {
			t.Fatal("expected profile, got nil")
		}
		if got.Profile.Difficulties[0] != "distraccion" {
			t.Errorf("expected difficulty %q, got %q", "distraccion", got.Profile.Difficulties[0])
		}
	})

	t.Run("rejects nil org_id", func(t *testing.T) {
		mock := &mocks.MockStudentProvider{}
		_, err := inclusion.NewGetStudentProfile(mock).Execute(ctx, inclusion.GetStudentProfileRequest{
			OrgID: uuid.Nil, StudentID: 1,
		})
		if !errors.Is(err, providers.ErrValidation) {
			t.Errorf("expected ErrValidation, got: %v", err)
		}
	})

	t.Run("rejects zero student_id", func(t *testing.T) {
		mock := &mocks.MockStudentProvider{}
		_, err := inclusion.NewGetStudentProfile(mock).Execute(ctx, inclusion.GetStudentProfileRequest{
			OrgID: testutil.TestOrgID, StudentID: 0,
		})
		if !errors.Is(err, providers.ErrValidation) {
			t.Errorf("expected ErrValidation, got: %v", err)
		}
	})

	t.Run("returns not found", func(t *testing.T) {
		mock := &mocks.MockStudentProvider{
			GetStudentFn: func(_ context.Context, _ uuid.UUID, _ int64) (*entities.Student, error) {
				return nil, errStudentNotFound
			},
		}
		_, err := inclusion.NewGetStudentProfile(mock).Execute(ctx, inclusion.GetStudentProfileRequest{
			OrgID: testutil.TestOrgID, StudentID: 999,
		})
		if !errors.Is(err, providers.ErrNotFound) {
			t.Errorf("expected ErrNotFound, got: %v", err)
		}
	})
}
