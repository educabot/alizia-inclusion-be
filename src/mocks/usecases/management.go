package mocks

import (
	"context"

	"github.com/stretchr/testify/mock"

	"github.com/educabot/alizia-inclusion-be/src/core/entities"
	mgmtuc "github.com/educabot/alizia-inclusion-be/src/core/usecases/management"
)

type MockListClassrooms struct {
	mock.Mock
}

func (m *MockListClassrooms) Execute(ctx context.Context, req mgmtuc.ListClassroomsRequest) ([]entities.Classroom, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]entities.Classroom), args.Error(1)
}

type MockGetClassroom struct {
	mock.Mock
}

func (m *MockGetClassroom) Execute(ctx context.Context, req mgmtuc.GetClassroomRequest) (*entities.Classroom, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.Classroom), args.Error(1)
}

type MockCreateClassroom struct {
	mock.Mock
}

func (m *MockCreateClassroom) Execute(ctx context.Context, req mgmtuc.CreateClassroomRequest) (*entities.Classroom, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.Classroom), args.Error(1)
}

type MockUpdateClassroom struct {
	mock.Mock
}

func (m *MockUpdateClassroom) Execute(ctx context.Context, req mgmtuc.UpdateClassroomRequest) (*entities.Classroom, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.Classroom), args.Error(1)
}

type MockDeleteClassroom struct {
	mock.Mock
}

func (m *MockDeleteClassroom) Execute(ctx context.Context, req mgmtuc.DeleteClassroomRequest) error {
	return m.Called(ctx, req).Error(0)
}

type MockListTeachers struct {
	mock.Mock
}

func (m *MockListTeachers) Execute(ctx context.Context, req mgmtuc.ListTeachersRequest) ([]entities.User, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]entities.User), args.Error(1)
}
