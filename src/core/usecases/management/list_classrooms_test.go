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

func TestListClassrooms_ReturnsClassrooms(t *testing.T) {
	ctx := context.Background()
	classrooms := []entities.Classroom{
		testutil.NewClassroom(1, "1A"),
		testutil.NewClassroom(2, "2B"),
	}
	mock := &mocks.MockClassroomProvider{
		ListFn: func(ctx context.Context, orgID uuid.UUID) ([]entities.Classroom, error) {
			return classrooms, nil
		},
	}

	uc := management.NewListClassrooms(mock)
	got, err := uc.Execute(ctx, management.ListClassroomsRequest{OrgID: testutil.TestOrgID})

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(got) != len(classrooms) {
		t.Fatalf("expected %d classrooms, got %d", len(classrooms), len(got))
	}
	if got[0].ID != classrooms[0].ID {
		t.Errorf("expected first classroom ID %d, got %d", classrooms[0].ID, got[0].ID)
	}
}

func TestListClassrooms_RejectsNilOrgID(t *testing.T) {
	ctx := context.Background()
	mock := &mocks.MockClassroomProvider{}

	uc := management.NewListClassrooms(mock)
	_, err := uc.Execute(ctx, management.ListClassroomsRequest{OrgID: uuid.Nil})

	if err == nil {
		t.Fatal("expected validation error, got nil")
	}
	if !errors.Is(err, providers.ErrValidation) {
		t.Errorf("expected ErrValidation, got %v", err)
	}
}
