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

func TestUpdateClassroom(t *testing.T) {
	ctx := context.Background()

	t.Run("updates classroom", func(t *testing.T) {
		newName := "Updated Name"
		newGrade := "4"
		newSection := "C"

		var updatedClassroom *entities.Classroom
		mock := &mocks.MockClassroomProvider{
			GetFn: func(ctx context.Context, orgID uuid.UUID, id int64) (*entities.Classroom, error) {
				c := testutil.NewClassroom(id, "Old Name")
				return &c, nil
			},
			UpdateFn: func(ctx context.Context, classroom *entities.Classroom) error {
				updatedClassroom = classroom
				return nil
			},
		}

		uc := management.NewUpdateClassroom(mock)
		got, err := uc.Execute(ctx, management.UpdateClassroomRequest{
			OrgID:       testutil.TestOrgID,
			ClassroomID: 1,
			Name:        &newName,
			Grade:       &newGrade,
			Section:     &newSection,
		})
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if updatedClassroom == nil {
			t.Fatal("expected UpdateFn to be called, got nil")
		}
		if updatedClassroom.Name != newName {
			t.Errorf("expected updated name %q, got %q", newName, updatedClassroom.Name)
		}
		if updatedClassroom.Grade == nil || *updatedClassroom.Grade != newGrade {
			t.Errorf("expected updated grade %q, got %v", newGrade, updatedClassroom.Grade)
		}
		if updatedClassroom.Section == nil || *updatedClassroom.Section != newSection {
			t.Errorf("expected updated section %q, got %v", newSection, updatedClassroom.Section)
		}
		if got.Name != newName {
			t.Errorf("expected returned name %q, got %q", newName, got.Name)
		}
	})

	t.Run("rejects nil org_id", func(t *testing.T) {
		mock := &mocks.MockClassroomProvider{}

		uc := management.NewUpdateClassroom(mock)
		_, err := uc.Execute(ctx, management.UpdateClassroomRequest{
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

		uc := management.NewUpdateClassroom(mock)
		_, err := uc.Execute(ctx, management.UpdateClassroomRequest{
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

		uc := management.NewUpdateClassroom(mock)
		newName := "New Name"
		_, err := uc.Execute(ctx, management.UpdateClassroomRequest{
			OrgID:       testutil.TestOrgID,
			ClassroomID: 99,
			Name:        &newName,
		})
		if err == nil {
			t.Fatal("expected not found error, got nil")
		}
		if !errors.Is(err, providers.ErrNotFound) {
			t.Errorf("expected ErrNotFound, got %v", err)
		}
	})
}
