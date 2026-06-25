package mocks

import (
	"context"

	"github.com/stretchr/testify/mock"

	"github.com/educabot/alizia-inclusion-be/src/core/providers"
)

type MockAIClient struct {
	mock.Mock
}

// Model returns a fixed identifier. It bypasses m.Called so every test that
// exercises an AI-backed usecase need not register an expectation for it.
func (m *MockAIClient) Model() string { return "mock-model" }

func (m *MockAIClient) Generate(ctx context.Context, params providers.GenerateParams) (string, error) {
	args := m.Called(ctx, params)
	return args.String(0), args.Error(1)
}

func (m *MockAIClient) Chat(ctx context.Context, messages []providers.ChatMessage) (*providers.ChatResponse, error) {
	args := m.Called(ctx, messages)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*providers.ChatResponse), args.Error(1)
}

func (m *MockAIClient) ChatWithTools(ctx context.Context, messages []providers.ChatMessage, tools []providers.ToolDefinition) (*providers.ChatResponse, error) {
	args := m.Called(ctx, messages, tools)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*providers.ChatResponse), args.Error(1)
}
