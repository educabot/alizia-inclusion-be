package ai

import (
	"context"
	"fmt"

	"github.com/educabot/alizia-inclusion-be/src/core/providers"
)

type StubClient struct{}

func NewStubClient() providers.AIClient {
	return &StubClient{}
}

func (s *StubClient) Generate(_ context.Context, params providers.GenerateParams) (string, error) {
	return fmt.Sprintf("[stub] generated content for: %s", params.UserPrompt), nil
}

func (s *StubClient) Chat(_ context.Context, _ []providers.ChatMessage) (*providers.ChatResponse, error) {
	return &providers.ChatResponse{Content: "[stub] chat response"}, nil
}

func (s *StubClient) ChatWithTools(_ context.Context, _ []providers.ChatMessage, _ []providers.ToolDefinition) (*providers.ChatResponse, error) {
	return &providers.ChatResponse{Content: "[stub] chat response"}, nil
}
