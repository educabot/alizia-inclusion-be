package providers

import "context"

type ChatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
	// ToolCallID links a tool-result message (role "tool") back to the assistant
	// tool call it answers. Empty for normal user/assistant/system messages.
	ToolCallID string `json:"tool_call_id,omitempty"`
	// ToolCalls echoes the tool calls an assistant message requested, so the
	// model can be reminded of them on the next turn of an agentic loop.
	ToolCalls []ToolCall `json:"tool_calls,omitempty"`
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
	// Model devuelve el identificador del modelo/deployment activo, para la traza
	// por turno (HU-6, T-6.5).
	Model() string
}
