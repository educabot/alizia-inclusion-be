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

func TestListClassroomStudents_ReturnsStudents(t *testing.T) {
	ctx := context.Background()
	expected := []entities.Student{
		testutil.NewStudent(1, 5, "Ana"),
		testutil.NewStudent(2, 5, "Lucas"),
	}
	mock := &mocks.MockStudentProvider{
		ListByClassroomFn: func(_ context.Context, _ uuid.UUID, classroomID int64) ([]entities.Student, error) {
			if classroomID != 5 {
				t.Errorf("expected classroomID 5, got %d", classroomID)
			}
			return expected, nil
		},
	}

	got, err := inclusion.NewListClassroomStudents(mock).Execute(ctx, inclusion.ListClassroomStudentsRequest{
		OrgID:       testutil.TestOrgID,
		ClassroomID: 5,
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != 2 {
		t.Errorf("got %d students, want 2", len(got))
	}
}

func TestListClassroomStudents_RejectsNilOrgID(t *testing.T) {
	ctx := context.Background()
	mock := &mocks.MockStudentProvider{}

	_, err := inclusion.NewListClassroomStudents(mock).Execute(ctx, inclusion.ListClassroomStudentsRequest{
		OrgID: uuid.Nil, ClassroomID: 1,
	})

	if !errors.Is(err, providers.ErrValidation) {
		t.Errorf("expected ErrValidation, got: %v", err)
	}
}

func TestListClassroomStudents_RejectsZeroClassroomID(t *testing.T) {
	ctx := context.Background()
	mock := &mocks.MockStudentProvider{}

	_, err := inclusion.NewListClassroomStudents(mock).Execute(ctx, inclusion.ListClassroomStudentsRequest{
		OrgID: testutil.TestOrgID, ClassroomID: 0,
	})

	if !errors.Is(err, providers.ErrValidation) {
		t.Errorf("expected ErrValidation, got: %v", err)
	}
}
