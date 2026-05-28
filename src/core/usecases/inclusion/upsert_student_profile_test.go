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

func TestUpsertStudentProfile_UpsertsProfile(t *testing.T) {
	ctx := context.Background()
	existing := testutil.NewStudent(1, 1, "Lucas")
	studentMock := &mocks.MockStudentProvider{
		GetStudentFn: func(_ context.Context, _ uuid.UUID, id int64) (*entities.Student, error) {
			s := existing
			return &s, nil
		},
	}
	var capturedProfile *entities.StudentProfile
	profileMock := &mocks.MockStudentProfileProvider{
		UpsertFn: func(_ context.Context, p *entities.StudentProfile) error {
			capturedProfile = p
			p.ID = 1
			return nil
		},
		GetByStudentIDFn: func(_ context.Context, studentID int64) (*entities.StudentProfile, error) {
			return &entities.StudentProfile{
				ID:           1,
				StudentID:    studentID,
				IsTransitory: true,
				Difficulties: []string{"distraccion", "motricidad_fina"},
			}, nil
		},
	}

	got, err := inclusion.NewUpsertStudentProfile(studentMock, profileMock).Execute(ctx, inclusion.UpsertStudentProfileRequest{
		OrgID:        testutil.TestOrgID,
		StudentID:    1,
		IsTransitory: true,
		Difficulties: []string{"distraccion", "motricidad_fina"},
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if capturedProfile == nil {
		t.Fatal("expected Upsert to be called")
	}
	if capturedProfile.StudentID != 1 {
		t.Errorf("expected student_id 1, got %d", capturedProfile.StudentID)
	}
	if !capturedProfile.IsTransitory {
		t.Error("expected is_transitory true")
	}
	if len(capturedProfile.Difficulties) != 2 {
		t.Errorf("expected 2 difficulties, got %d", len(capturedProfile.Difficulties))
	}
	if got == nil {
		t.Fatal("expected non-nil result")
	}
}

func TestUpsertStudentProfile_RejectsNilOrgID(t *testing.T) {
	ctx := context.Background()
	studentMock := &mocks.MockStudentProvider{}
	profileMock := &mocks.MockStudentProfileProvider{}

	_, err := inclusion.NewUpsertStudentProfile(studentMock, profileMock).Execute(ctx, inclusion.UpsertStudentProfileRequest{
		OrgID: uuid.Nil, StudentID: 1, Difficulties: []string{},
	})

	if !errors.Is(err, providers.ErrValidation) {
		t.Errorf("expected ErrValidation, got: %v", err)
	}
}

func TestUpsertStudentProfile_RejectsZeroStudentID(t *testing.T) {
	ctx := context.Background()
	studentMock := &mocks.MockStudentProvider{}
	profileMock := &mocks.MockStudentProfileProvider{}

	_, err := inclusion.NewUpsertStudentProfile(studentMock, profileMock).Execute(ctx, inclusion.UpsertStudentProfileRequest{
		OrgID: testutil.TestOrgID, StudentID: 0, Difficulties: []string{},
	})

	if !errors.Is(err, providers.ErrValidation) {
		t.Errorf("expected ErrValidation, got: %v", err)
	}
}

func TestUpsertStudentProfile_ReturnsErrorIfStudentNotFound(t *testing.T) {
	ctx := context.Background()
	studentMock := &mocks.MockStudentProvider{
		GetStudentFn: func(_ context.Context, _ uuid.UUID, _ int64) (*entities.Student, error) {
			return nil, errStudentNotFound
		},
	}
	profileMock := &mocks.MockStudentProfileProvider{}

	_, err := inclusion.NewUpsertStudentProfile(studentMock, profileMock).Execute(ctx, inclusion.UpsertStudentProfileRequest{
		OrgID: testutil.TestOrgID, StudentID: 999, Difficulties: []string{},
	})

	if !errors.Is(err, providers.ErrNotFound) {
		t.Errorf("expected ErrNotFound, got: %v", err)
	}
}
