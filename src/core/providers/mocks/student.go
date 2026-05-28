package mocks

import (
	"context"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"

	"github.com/educabot/alizia-inclusion-be/src/core/entities"
)

type MockStudentProvider struct {
	mock.Mock
}

func (m *MockStudentProvider) GetStudent(ctx context.Context, orgID uuid.UUID, id int64) (*entities.Student, error) {
	args := m.Called(ctx, orgID, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.Student), args.Error(1)
}

func (m *MockStudentProvider) ListByClassroom(ctx context.Context, orgID uuid.UUID, classroomID int64) ([]entities.Student, error) {
	args := m.Called(ctx, orgID, classroomID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]entities.Student), args.Error(1)
}

func (m *MockStudentProvider) List(ctx context.Context, orgID uuid.UUID) ([]entities.Student, error) {
	args := m.Called(ctx, orgID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]entities.Student), args.Error(1)
}

func (m *MockStudentProvider) Create(ctx context.Context, student *entities.Student) error {
	args := m.Called(ctx, student)
	return args.Error(0)
}

func (m *MockStudentProvider) Update(ctx context.Context, student *entities.Student) error {
	args := m.Called(ctx, student)
	return args.Error(0)
}

func (m *MockStudentProvider) Delete(ctx context.Context, orgID uuid.UUID, id int64) error {
	args := m.Called(ctx, orgID, id)
	return args.Error(0)
}
