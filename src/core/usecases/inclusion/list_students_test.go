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

func TestListStudents(t *testing.T) {
	ctx := context.Background()

	t.Run("returns all students", func(t *testing.T) {
		expected := []entities.Student{
			testutil.NewStudent(1, 1, "Lucas"),
			testutil.NewStudent(2, 1, "Ana"),
		}
		mock := &mocks.MockStudentProvider{
			ListFn: func(_ context.Context, _ uuid.UUID) ([]entities.Student, error) {
				return expected, nil
			},
		}

		got, err := inclusion.NewListStudents(mock).Execute(ctx, inclusion.ListStudentsRequest{
			OrgID: testutil.TestOrgID,
		})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(got) != 2 {
			t.Errorf("got %d students, want 2", len(got))
		}
	})

	t.Run("filters by classroom", func(t *testing.T) {
		var capturedClassroomID int64
		mock := &mocks.MockStudentProvider{
			ListByClassroomFn: func(_ context.Context, _ uuid.UUID, classroomID int64) ([]entities.Student, error) {
				capturedClassroomID = classroomID
				return []entities.Student{testutil.NewStudent(1, classroomID, "Lucas")}, nil
			},
		}

		classroomID := int64(5)
		got, err := inclusion.NewListStudents(mock).Execute(ctx, inclusion.ListStudentsRequest{
			OrgID:       testutil.TestOrgID,
			ClassroomID: &classroomID,
		})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if capturedClassroomID != 5 {
			t.Errorf("expected classroomID 5, got %d", capturedClassroomID)
		}
		if len(got) != 1 {
			t.Errorf("got %d students, want 1", len(got))
		}
	})

	t.Run("rejects nil org_id", func(t *testing.T) {
		mock := &mocks.MockStudentProvider{}
		_, err := inclusion.NewListStudents(mock).Execute(ctx, inclusion.ListStudentsRequest{
			OrgID: uuid.Nil,
		})
		if !errors.Is(err, providers.ErrValidation) {
			t.Errorf("expected ErrValidation, got: %v", err)
		}
	})
}
