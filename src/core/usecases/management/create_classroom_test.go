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

func TestCreateClassroom(t *testing.T) {
	ctx := context.Background()

	t.Run("creates classroom", func(t *testing.T) {
		grade := "3"
		section := "B"

		var capturedClassroom *entities.Classroom
		mock := &mocks.MockClassroomProvider{
			CreateFn: func(ctx context.Context, classroom *entities.Classroom) error {
				capturedClassroom = classroom
				classroom.ID = 42
				return nil
			},
		}

		uc := management.NewCreateClassroom(mock)
		got, err := uc.Execute(ctx, management.CreateClassroomRequest{
			OrgID:   testutil.TestOrgID,
			Name:    "3B Math",
			Grade:   &grade,
			Section: &section,
		})
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if capturedClassroom == nil {
			t.Fatal("expected CreateFn to be called with a classroom, got nil")
		}
		if capturedClassroom.OrganizationID != testutil.TestOrgID {
			t.Errorf("expected OrgID %v, got %v", testutil.TestOrgID, capturedClassroom.OrganizationID)
		}
		if capturedClassroom.Name != "3B Math" {
			t.Errorf("expected name %q, got %q", "3B Math", capturedClassroom.Name)
		}
		if capturedClassroom.Grade == nil || *capturedClassroom.Grade != grade {
			t.Errorf("expected grade %q, got %v", grade, capturedClassroom.Grade)
		}
		if capturedClassroom.Section == nil || *capturedClassroom.Section != section {
			t.Errorf("expected section %q, got %v", section, capturedClassroom.Section)
		}
		if got.ID != 42 {
			t.Errorf("expected returned classroom ID 42, got %d", got.ID)
		}
	})

	t.Run("rejects nil org_id", func(t *testing.T) {
		mock := &mocks.MockClassroomProvider{}

		uc := management.NewCreateClassroom(mock)
		_, err := uc.Execute(ctx, management.CreateClassroomRequest{
			OrgID: uuid.Nil,
			Name:  "1A",
		})
		if err == nil {
			t.Fatal("expected validation error, got nil")
		}
		if !errors.Is(err, providers.ErrValidation) {
			t.Errorf("expected ErrValidation, got %v", err)
		}
	})

	t.Run("rejects empty name", func(t *testing.T) {
		mock := &mocks.MockClassroomProvider{}

		uc := management.NewCreateClassroom(mock)
		_, err := uc.Execute(ctx, management.CreateClassroomRequest{
			OrgID: testutil.TestOrgID,
			Name:  "",
		})
		if err == nil {
			t.Fatal("expected validation error, got nil")
		}
		if !errors.Is(err, providers.ErrValidation) {
			t.Errorf("expected ErrValidation, got %v", err)
		}
	})

	t.Run("propagates provider error", func(t *testing.T) {
		mock := &mocks.MockClassroomProvider{
			CreateFn: func(ctx context.Context, classroom *entities.Classroom) error {
				return errDB
			},
		}

		uc := management.NewCreateClassroom(mock)
		_, err := uc.Execute(ctx, management.CreateClassroomRequest{
			OrgID: testutil.TestOrgID,
			Name:  "1A",
		})
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if !errors.Is(err, errDB) {
			t.Errorf("expected errDB, got %v", err)
		}
	})
}
