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

type azureMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type azureRequest struct {
	Messages    []azureMessage `json:"messages"`
	Temperature *float32       `json:"temperature,omitempty"`
	MaxTokens   *int           `json:"max_tokens,omitempty"`
}

type azureChoice struct {
	Message azureMessage `json:"message"`
}

type azureResponse struct {
	Choices []azureChoice `json:"choices"`
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

func (a *AzureClient) ChatWithTools(ctx context.Context, messages []providers.ChatMessage, _ []providers.ToolDefinition) (*providers.ChatResponse, error) {
	azMsgs := make([]azureMessage, 0, len(messages))
	for _, m := range messages {
		azMsgs = append(azMsgs, azureMessage{Role: m.Role, Content: m.Content})
	}

	resp, err := a.doRequest(ctx, azureRequest{Messages: azMsgs})
	if err != nil {
		return nil, err
	}

	if len(resp.Choices) == 0 {
		return &providers.ChatResponse{Content: ""}, nil
	}

	return &providers.ChatResponse{Content: resp.Choices[0].Message.Content}, nil
}
