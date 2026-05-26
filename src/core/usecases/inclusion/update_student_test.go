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

func TestUpdateStudent(t *testing.T) {
	ctx := context.Background()

	t.Run("updates student", func(t *testing.T) {
		existing := testutil.NewStudent(1, 1, "Old Name")
		mock := &mocks.MockStudentProvider{
			GetStudentFn: func(_ context.Context, _ uuid.UUID, id int64) (*entities.Student, error) {
				s := existing
				return &s, nil
			},
			UpdateFn: func(_ context.Context, s *entities.Student) error {
				return nil
			},
		}

		newName := "New Name"
		got, err := inclusion.NewUpdateStudent(mock).Execute(ctx, inclusion.UpdateStudentRequest{
			OrgID:     testutil.TestOrgID,
			StudentID: 1,
			Name:      &newName,
		})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if got.Name != newName {
			t.Errorf("expected name %q, got %q", newName, got.Name)
		}
	})

	t.Run("rejects nil org_id", func(t *testing.T) {
		mock := &mocks.MockStudentProvider{}
		_, err := inclusion.NewUpdateStudent(mock).Execute(ctx, inclusion.UpdateStudentRequest{
			OrgID: uuid.Nil, StudentID: 1,
		})
		if !errors.Is(err, providers.ErrValidation) {
			t.Errorf("expected ErrValidation, got: %v", err)
		}
	})

	t.Run("rejects zero student_id", func(t *testing.T) {
		mock := &mocks.MockStudentProvider{}
		_, err := inclusion.NewUpdateStudent(mock).Execute(ctx, inclusion.UpdateStudentRequest{
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
		_, err := inclusion.NewUpdateStudent(mock).Execute(ctx, inclusion.UpdateStudentRequest{
			OrgID: testutil.TestOrgID, StudentID: 999,
		})
		if !errors.Is(err, providers.ErrNotFound) {
			t.Errorf("expected ErrNotFound, got: %v", err)
		}
	})
}
