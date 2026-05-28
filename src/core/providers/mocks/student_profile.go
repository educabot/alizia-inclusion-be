package mocks

import (
	"context"

	"github.com/stretchr/testify/mock"

	"github.com/educabot/alizia-inclusion-be/src/core/entities"
)

type MockStudentProfileProvider struct {
	mock.Mock
}

func (m *MockStudentProfileProvider) GetByStudentID(ctx context.Context, studentID int64) (*entities.StudentProfile, error) {
	args := m.Called(ctx, studentID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.StudentProfile), args.Error(1)
}

func (m *MockStudentProfileProvider) Upsert(ctx context.Context, profile *entities.StudentProfile) error {
	args := m.Called(ctx, profile)
	return args.Error(0)
}
