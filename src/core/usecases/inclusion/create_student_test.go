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

func TestCreateStudent_CreatesStudent(t *testing.T) {
	ctx := context.Background()

	var captured *entities.Student
	mock := &mocks.MockStudentProvider{
		CreateFn: func(_ context.Context, s *entities.Student) error {
			captured = s
			s.ID = 10
			return nil
		},
	}

	got, err := inclusion.NewCreateStudent(mock).Execute(ctx, inclusion.CreateStudentRequest{
		OrgID:       testutil.TestOrgID,
		ClassroomID: 1,
		Name:        "Juan",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if captured.Name != "Juan" {
		t.Errorf("expected name %q, got %q", "Juan", captured.Name)
	}
	if captured.ClassroomID != 1 {
		t.Errorf("expected classroom_id 1, got %d", captured.ClassroomID)
	}
	if got.ID != 10 {
		t.Errorf("expected ID 10, got %d", got.ID)
	}
}

func TestCreateStudent_RejectsNilOrgID(t *testing.T) {
	ctx := context.Background()

	mock := &mocks.MockStudentProvider{}
	_, err := inclusion.NewCreateStudent(mock).Execute(ctx, inclusion.CreateStudentRequest{
		OrgID: uuid.Nil, ClassroomID: 1, Name: "X",
	})
	if !errors.Is(err, providers.ErrValidation) {
		t.Errorf("expected ErrValidation, got: %v", err)
	}
}

func TestCreateStudent_RejectsZeroClassroomID(t *testing.T) {
	ctx := context.Background()

	mock := &mocks.MockStudentProvider{}
	_, err := inclusion.NewCreateStudent(mock).Execute(ctx, inclusion.CreateStudentRequest{
		OrgID: testutil.TestOrgID, ClassroomID: 0, Name: "X",
	})
	if !errors.Is(err, providers.ErrValidation) {
		t.Errorf("expected ErrValidation, got: %v", err)
	}
}

func TestCreateStudent_RejectsEmptyName(t *testing.T) {
	ctx := context.Background()

	mock := &mocks.MockStudentProvider{}
	_, err := inclusion.NewCreateStudent(mock).Execute(ctx, inclusion.CreateStudentRequest{
		OrgID: testutil.TestOrgID, ClassroomID: 1, Name: "",
	})
	if !errors.Is(err, providers.ErrValidation) {
		t.Errorf("expected ErrValidation, got: %v", err)
	}
}
