package mocks

import (
	"context"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"

	"github.com/educabot/alizia-inclusion-be/src/core/entities"
)

type MockClassroomProvider struct {
	mock.Mock
}

func (m *MockClassroomProvider) List(ctx context.Context, orgID uuid.UUID) ([]entities.Classroom, error) {
	args := m.Called(ctx, orgID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]entities.Classroom), args.Error(1)
}

func (m *MockClassroomProvider) Get(ctx context.Context, orgID uuid.UUID, id int64) (*entities.Classroom, error) {
	args := m.Called(ctx, orgID, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.Classroom), args.Error(1)
}

func (m *MockClassroomProvider) Create(ctx context.Context, classroom *entities.Classroom) error {
	args := m.Called(ctx, classroom)
	return args.Error(0)
}

func (m *MockClassroomProvider) Update(ctx context.Context, classroom *entities.Classroom) error {
	args := m.Called(ctx, classroom)
	return args.Error(0)
}

func (m *MockClassroomProvider) Delete(ctx context.Context, orgID uuid.UUID, id int64) error {
	args := m.Called(ctx, orgID, id)
	return args.Error(0)
}
