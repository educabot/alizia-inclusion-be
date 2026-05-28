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

func TestUpdateStudent_UpdatesStudent(t *testing.T) {
	ctx := context.Background()
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
}

func TestUpdateStudent_RejectsNilOrgID(t *testing.T) {
	ctx := context.Background()
	mock := &mocks.MockStudentProvider{}

	_, err := inclusion.NewUpdateStudent(mock).Execute(ctx, inclusion.UpdateStudentRequest{
		OrgID: uuid.Nil, StudentID: 1,
	})

	if !errors.Is(err, providers.ErrValidation) {
		t.Errorf("expected ErrValidation, got: %v", err)
	}
}

func TestUpdateStudent_RejectsZeroStudentID(t *testing.T) {
	ctx := context.Background()
	mock := &mocks.MockStudentProvider{}

	_, err := inclusion.NewUpdateStudent(mock).Execute(ctx, inclusion.UpdateStudentRequest{
		OrgID: testutil.TestOrgID, StudentID: 0,
	})

	if !errors.Is(err, providers.ErrValidation) {
		t.Errorf("expected ErrValidation, got: %v", err)
	}
}

func TestUpdateStudent_ReturnsNotFound(t *testing.T) {
	ctx := context.Background()
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
}
