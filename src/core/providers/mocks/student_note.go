package mocks

import (
	"context"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"

	"github.com/educabot/alizia-inclusion-be/src/core/entities"
)

type MockStudentNoteProvider struct {
	mock.Mock
}

func (m *MockStudentNoteProvider) ListByStudent(ctx context.Context, orgID uuid.UUID, studentID, userID int64) ([]entities.StudentNote, error) {
	args := m.Called(ctx, orgID, studentID, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]entities.StudentNote), args.Error(1)
}

func (m *MockStudentNoteProvider) Create(ctx context.Context, note *entities.StudentNote) error {
	return m.Called(ctx, note).Error(0)
}
