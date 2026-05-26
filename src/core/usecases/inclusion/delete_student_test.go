package inclusion_test

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"

	"github.com/educabot/alizia-inclusion-be/src/core/providers"
	"github.com/educabot/alizia-inclusion-be/src/core/providers/mocks"
	"github.com/educabot/alizia-inclusion-be/src/core/usecases/inclusion"
	"github.com/educabot/alizia-inclusion-be/src/testutil"
)

func TestDeleteStudent(t *testing.T) {
	ctx := context.Background()

	t.Run("deletes student", func(t *testing.T) {
		var calledWith int64
		mock := &mocks.MockStudentProvider{
			DeleteFn: func(_ context.Context, _ uuid.UUID, id int64) error {
				calledWith = id
				return nil
			},
		}

		err := inclusion.NewDeleteStudent(mock).Execute(ctx, inclusion.DeleteStudentRequest{
			OrgID:     testutil.TestOrgID,
			StudentID: 5,
		})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if calledWith != 5 {
			t.Errorf("expected delete called with ID 5, got %d", calledWith)
		}
	})

	t.Run("rejects nil org_id", func(t *testing.T) {
		mock := &mocks.MockStudentProvider{}
		err := inclusion.NewDeleteStudent(mock).Execute(ctx, inclusion.DeleteStudentRequest{
			OrgID: uuid.Nil, StudentID: 1,
		})
		if !errors.Is(err, providers.ErrValidation) {
			t.Errorf("expected ErrValidation, got: %v", err)
		}
	})

	t.Run("rejects zero student_id", func(t *testing.T) {
		mock := &mocks.MockStudentProvider{}
		err := inclusion.NewDeleteStudent(mock).Execute(ctx, inclusion.DeleteStudentRequest{
			OrgID: testutil.TestOrgID, StudentID: 0,
		})
		if !errors.Is(err, providers.ErrValidation) {
			t.Errorf("expected ErrValidation, got: %v", err)
		}
	})
}
