package ai

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/educabot/alizia-inclusion-be/src/core/providers"
)

type AzureClient struct {
	endpoint   string
	apiKey     string
	deployment string
	httpClient *http.Client
}

func NewAzureClient(endpoint, apiKey, deployment string) providers.AIClient {
	endpoint = strings.TrimRight(endpoint, "/")
	return &AzureClient{
		endpoint:   endpoint,
		apiKey:     apiKey,
		deployment: deployment,
		httpClient: &http.Client{},
	}
}

// Model returns the Azure deployment name as the model identifier.
func (c *AzureClient) Model() string { return c.deployment }

type azureMessage struct {
	Role       string          `json:"role"`
	Content    string          `json:"content"`
	ToolCallID string          `json:"tool_call_id,omitempty"`
	ToolCalls  []azureToolCall `json:"tool_calls,omitempty"`
}

type azureToolFunction struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Parameters  any    `json:"parameters,omitempty"`
}

type azureTool struct {
	Type     string            `json:"type"`
	Function azureToolFunction `json:"function"`
}

type azureRequest struct {
	Messages    []azureMessage `json:"messages"`
	Temperature *float32       `json:"temperature,omitempty"`
	MaxTokens   *int           `json:"max_tokens,omitempty"`
	Tools       []azureTool    `json:"tools,omitempty"`
}

type azureToolCall struct {
	ID       string `json:"id"`
	Type     string `json:"type"`
	Function struct {
		Name      string `json:"name"`
		Arguments string `json:"arguments"`
	} `json:"function"`
}

type azureRespMessage struct {
	Role      string          `json:"role"`
	Content   string          `json:"content"`
	ToolCalls []azureToolCall `json:"tool_calls,omitempty"`
}

type azureChoice struct {
	Message azureRespMessage `json:"message"`
}

type azureUsage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

type azureResponse struct {
	Choices []azureChoice `json:"choices"`
	Usage   *azureUsage   `json:"usage,omitempty"`
	Error   *azureError   `json:"error,omitempty"`
}

type azureError struct {
	Message string `json:"message"`
	Code    string `json:"code"`
}

func (a *AzureClient) doRequest(ctx context.Context, req azureRequest) (*azureResponse, error) {
	body, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("marshal request: %w", err)
	}

	url := fmt.Sprintf("%s/openai/deployments/%s/chat/completions?api-version=2024-10-21", a.endpoint, a.deployment)
	httpReq, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("api-key", a.apiKey)

	resp, err := a.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("do request: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("azure openai error (status %d): %s", resp.StatusCode, string(respBody))
	}

	var azResp azureResponse
	if err := json.Unmarshal(respBody, &azResp); err != nil {
		return nil, fmt.Errorf("unmarshal response: %w", err)
	}

	if azResp.Error != nil {
		return nil, fmt.Errorf("azure openai: %s (code: %s)", azResp.Error.Message, azResp.Error.Code)
	}

	return &azResp, nil
}

func (a *AzureClient) Generate(ctx context.Context, params providers.GenerateParams) (string, error) {
	messages := []azureMessage{
		{Role: "system", Content: params.SystemPrompt},
		{Role: "user", Content: params.UserPrompt},
	}

	resp, err := a.doRequest(ctx, azureRequest{
		Messages:    messages,
		Temperature: params.Temperature,
		MaxTokens:   params.MaxTokens,
	})
	if err != nil {
		return "", err
	}

	if len(resp.Choices) == 0 {
		return "", fmt.Errorf("azure generate: empty response")
	}

	return resp.Choices[0].Message.Content, nil
}

func (a *AzureClient) Chat(ctx context.Context, messages []providers.ChatMessage) (*providers.ChatResponse, error) {
	return a.ChatWithTools(ctx, messages, nil)
}

func (a *AzureClient) ChatWithTools(ctx context.Context, messages []providers.ChatMessage, tools []providers.ToolDefinition) (*providers.ChatResponse, error) {
	azMsgs := make([]azureMessage, 0, len(messages))
	for _, m := range messages {
		am := azureMessage{Role: m.Role, Content: m.Content, ToolCallID: m.ToolCallID}
		for _, tc := range m.ToolCalls {
			ac := azureToolCall{ID: tc.ID, Type: "function"}
			ac.Function.Name = tc.Name
			ac.Function.Arguments = tc.Arguments
			am.ToolCalls = append(am.ToolCalls, ac)
		}
		azMsgs = append(azMsgs, am)
	}

	var azTools []azureTool
	if len(tools) > 0 {
		azTools = make([]azureTool, 0, len(tools))
		for _, t := range tools {
			azTools = append(azTools, azureTool{
				Type: "function",
				Function: azureToolFunction{
					Name:        t.Name,
					Description: t.Description,
					Parameters:  t.Parameters,
				},
			})
		}
	}

	resp, err := a.doRequest(ctx, azureRequest{Messages: azMsgs, Tools: azTools})
	if err != nil {
		return nil, err
	}

	if len(resp.Choices) == 0 {
		return &providers.ChatResponse{Content: ""}, nil
	}

	msg := resp.Choices[0].Message
	out := &providers.ChatResponse{Content: msg.Content}
	for _, tc := range msg.ToolCalls {
		out.ToolCalls = append(out.ToolCalls, providers.ToolCall{
			ID:        tc.ID,
			Name:      tc.Function.Name,
			Arguments: tc.Function.Arguments,
		})
	}
	if resp.Usage != nil {
		out.Usage = &providers.TokenUsage{
			PromptTokens:     resp.Usage.PromptTokens,
			CompletionTokens: resp.Usage.CompletionTokens,
			TotalTokens:      resp.Usage.TotalTokens,
		}
	}
	return out, nil
}
