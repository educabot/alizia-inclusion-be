package providers

import "context"

type ChatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ToolDefinition struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Parameters  any    `json:"parameters"`
}

type ToolCall struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Arguments string `json:"arguments"`
}

type TokenUsage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

type ChatResponse struct {
	Content   string      `json:"content"`
	ToolCalls []ToolCall  `json:"tool_calls,omitempty"`
	Usage     *TokenUsage `json:"usage,omitempty"`
}

type GenerateParams struct {
	SystemPrompt string
	UserPrompt   string
	Model        string
	Temperature  *float32
	MaxTokens    *int
	JSONResponse bool
}

type AIClient interface {
	Generate(ctx context.Context, params GenerateParams) (string, error)
	Chat(ctx context.Context, messages []ChatMessage) (*ChatResponse, error)
	ChatWithTools(ctx context.Context, messages []ChatMessage, tools []ToolDefinition) (*ChatResponse, error)
}
