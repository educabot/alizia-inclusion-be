package mocks

import (
	"context"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"

	"github.com/educabot/alizia-inclusion-be/src/core/entities"
)

type MockTeacherProfileProvider struct {
	mock.Mock
}

func (m *MockTeacherProfileProvider) GetByUserID(ctx context.Context, orgID uuid.UUID, userID int64) (*entities.TeacherProfile, error) {
	args := m.Called(ctx, orgID, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.TeacherProfile), args.Error(1)
}

type MockSituationCatalogProvider struct {
	mock.Mock
}

func (m *MockSituationCatalogProvider) List(ctx context.Context, orgID uuid.UUID) ([]entities.Situation, error) {
	args := m.Called(ctx, orgID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]entities.Situation), args.Error(1)
}

type MockDiagnosisProvider struct {
	mock.Mock
}

func (m *MockDiagnosisProvider) ListByStudentProfile(ctx context.Context, studentProfileID int64) ([]entities.StudentDiagnosis, error) {
	args := m.Called(ctx, studentProfileID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]entities.StudentDiagnosis), args.Error(1)
}

type MockPPIProvider struct {
	mock.Mock
}

func (m *MockPPIProvider) GetByStudentID(ctx context.Context, orgID uuid.UUID, studentID int64) (*entities.PPI, error) {
	args := m.Called(ctx, orgID, studentID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.PPI), args.Error(1)
}
