package mocks

import (
	"context"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"

	"github.com/educabot/alizia-inclusion-be/src/core/entities"
	"github.com/educabot/alizia-inclusion-be/src/core/providers"
)

type MockMessageFeedbackProvider struct {
	mock.Mock
}

func (m *MockMessageFeedbackProvider) MessageContext(ctx context.Context, messageID int64) (int64, uuid.UUID, error) {
	args := m.Called(ctx, messageID)
	return args.Get(0).(int64), args.Get(1).(uuid.UUID), args.Error(2)
}

func (m *MockMessageFeedbackProvider) Upsert(ctx context.Context, fb *entities.MessageFeedback) error {
	return m.Called(ctx, fb).Error(0)
}

func (m *MockMessageFeedbackProvider) Delete(ctx context.Context, orgID uuid.UUID, messageID, userID int64) error {
	return m.Called(ctx, orgID, messageID, userID).Error(0)
}

func (m *MockMessageFeedbackProvider) List(ctx context.Context, orgID uuid.UUID, rating string) ([]providers.MessageFeedbackReview, error) {
	args := m.Called(ctx, orgID, rating)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]providers.MessageFeedbackReview), args.Error(1)
}
