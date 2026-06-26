package ai_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	ai "github.com/educabot/alizia-inclusion-be/src/repositories/ai"

	"github.com/educabot/alizia-inclusion-be/src/core/providers"
)

func validAzureResponse(content string) map[string]any {
	return map[string]any{
		"choices": []map[string]any{
			{"message": map[string]string{"role": "assistant", "content": content}},
		},
	}
}

func writeJSON(t *testing.T, w http.ResponseWriter, v any) {
	t.Helper()
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(v); err != nil {
		t.Errorf("writeJSON: %v", err)
	}
}

func TestAzureClient_Generate_ReturnsGeneratedContent(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "test-key", r.Header.Get("api-key"))
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))
		writeJSON(t, w, validAzureResponse("test response"))
	}))
	defer server.Close()

	client := ai.NewAzureClient(server.URL, "test-key", "gpt-4o-mini", "2024-12-01-preview")
	got, err := client.Generate(context.Background(), providers.GenerateParams{
		SystemPrompt: "you are helpful",
		UserPrompt:   "hello",
	})

	require.NoError(t, err)
	assert.Equal(t, "test response", got)
}

func TestAzureClient_Generate_ReturnsErrorForNon200Status(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}))
	defer server.Close()

	client := ai.NewAzureClient(server.URL, "test-key", "gpt-4o-mini", "2024-12-01-preview")
	_, err := client.Generate(context.Background(), providers.GenerateParams{
		SystemPrompt: "you are helpful",
		UserPrompt:   "hello",
	})

	require.Error(t, err)
	assert.Contains(t, err.Error(), "500")
}

func TestAzureClient_Generate_ReturnsErrorForEmptyChoices(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		writeJSON(t, w, map[string]any{"choices": []any{}})
	}))
	defer server.Close()

	client := ai.NewAzureClient(server.URL, "test-key", "gpt-4o-mini", "2024-12-01-preview")
	_, err := client.Generate(context.Background(), providers.GenerateParams{
		SystemPrompt: "you are helpful",
		UserPrompt:   "hello",
	})

	require.Error(t, err)
	assert.Contains(t, err.Error(), "empty response")
}

func TestAzureClient_Chat_ReturnsChatResponse(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		writeJSON(t, w, validAzureResponse("chat reply"))
	}))
	defer server.Close()

	client := ai.NewAzureClient(server.URL, "test-key", "gpt-4o-mini", "2024-12-01-preview")
	resp, err := client.Chat(context.Background(), []providers.ChatMessage{
		{Role: "user", Content: "hello"},
	})

	require.NoError(t, err)
	require.NotNil(t, resp)
	assert.Equal(t, "chat reply", resp.Content)
}

func TestAzureClient_Chat_SendsMessagesCorrectly(t *testing.T) {
	var capturedBody map[string]any

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		err := json.NewDecoder(r.Body).Decode(&capturedBody)
		require.NoError(t, err)
		writeJSON(t, w, validAzureResponse("ok"))
	}))
	defer server.Close()

	client := ai.NewAzureClient(server.URL, "test-key", "gpt-4o-mini", "2024-12-01-preview")
	messages := []providers.ChatMessage{
		{Role: "system", Content: "be helpful"},
		{Role: "user", Content: "what is Go?"},
	}
	_, err := client.Chat(context.Background(), messages)
	require.NoError(t, err)

	rawMessages, ok := capturedBody["messages"].([]any)
	require.True(t, ok, "messages field missing or not a slice in request body")
	require.Len(t, rawMessages, 2)

	first, ok := rawMessages[0].(map[string]any)
	require.True(t, ok, "first message is not a map")
	assert.Equal(t, "system", first["role"])
	assert.Equal(t, "be helpful", first["content"])

	second, ok := rawMessages[1].(map[string]any)
	require.True(t, ok, "second message is not a map")
	assert.Equal(t, "user", second["role"])
	assert.Equal(t, "what is Go?", second["content"])
}

func TestAzureClient_ChatWithTools_ReturnsResponse(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		writeJSON(t, w, validAzureResponse("tool reply"))
	}))
	defer server.Close()

	client := ai.NewAzureClient(server.URL, "test-key", "gpt-4o-mini", "2024-12-01-preview")
	resp, err := client.ChatWithTools(
		context.Background(),
		[]providers.ChatMessage{{Role: "user", Content: "use a tool"}},
		[]providers.ToolDefinition{{Name: "my_tool", Description: "does stuff"}},
	)

	require.NoError(t, err)
	require.NotNil(t, resp)
	assert.Equal(t, "tool reply", resp.Content)
}

func TestAzureClient_ChatWithTools_SendsToolDefinitionsInOpenAIFunctionFormat(t *testing.T) {
	var capturedBody map[string]any
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		err := json.NewDecoder(r.Body).Decode(&capturedBody)
		require.NoError(t, err)
		writeJSON(t, w, validAzureResponse("ok"))
	}))
	defer server.Close()

	client := ai.NewAzureClient(server.URL, "test-key", "gpt-4o-mini", "2024-12-01-preview")
	_, err := client.ChatWithTools(
		context.Background(),
		[]providers.ChatMessage{{Role: "user", Content: "buscar dispositivos"}},
		[]providers.ToolDefinition{{
			Name:        "search_devices",
			Description: "busca dispositivos",
			Parameters:  map[string]any{"type": "object"},
		}},
	)
	require.NoError(t, err)

	tools, ok := capturedBody["tools"].([]any)
	require.True(t, ok)
	require.Len(t, tools, 1)

	tool, ok := tools[0].(map[string]any)
	require.True(t, ok, "tool is not a map")
	assert.Equal(t, "function", tool["type"])

	fn, ok := tool["function"].(map[string]any)
	require.True(t, ok, "tool.function is not a map")
	assert.Equal(t, "search_devices", fn["name"])
}

func TestAzureClient_ChatWithTools_ParsesToolCallsFromResponse(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		writeJSON(t, w, map[string]any{
			"choices": []map[string]any{
				{"message": map[string]any{
					"role":    "assistant",
					"content": "",
					"tool_calls": []map[string]any{
						{
							"id":   "call_1",
							"type": "function",
							"function": map[string]any{
								"name":      "search_devices",
								"arguments": `{"query":"timer"}`,
							},
						},
					},
				}},
			},
		})
	}))
	defer server.Close()

	client := ai.NewAzureClient(server.URL, "test-key", "gpt-4o-mini", "2024-12-01-preview")
	resp, err := client.ChatWithTools(
		context.Background(),
		[]providers.ChatMessage{{Role: "user", Content: "buscar timer"}},
		[]providers.ToolDefinition{{Name: "search_devices", Description: "busca"}},
	)
	require.NoError(t, err)
	require.Len(t, resp.ToolCalls, 1)

	tc := resp.ToolCalls[0]
	assert.Equal(t, "call_1", tc.ID)
	assert.Equal(t, "search_devices", tc.Name)
	assert.Equal(t, `{"query":"timer"}`, tc.Arguments)
}

func TestAzureClient_ChatWithTools_SerializesAssistantToolCallsAndToolResultMessagesInOpenAIFormat(t *testing.T) {
	var capturedBody map[string]any
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		err := json.NewDecoder(r.Body).Decode(&capturedBody)
		require.NoError(t, err)
		writeJSON(t, w, validAzureResponse("ok"))
	}))
	defer server.Close()

	client := ai.NewAzureClient(server.URL, "test-key", "gpt-4o-mini", "2024-12-01-preview")
	_, err := client.ChatWithTools(
		context.Background(),
		[]providers.ChatMessage{
			{Role: "user", Content: "que dispositivo uso?"},
			{Role: "assistant", ToolCalls: []providers.ToolCall{
				{ID: "call_9", Name: "list_devices", Arguments: "{}"},
			}},
			{Role: "tool", ToolCallID: "call_9", Content: `{"devices":[]}`},
		},
		[]providers.ToolDefinition{{Name: "list_devices", Description: "lista"}},
	)
	require.NoError(t, err)

	msgs, ok := capturedBody["messages"].([]any)
	require.True(t, ok)
	require.Len(t, msgs, 3)

	assistant, ok := msgs[1].(map[string]any)
	require.True(t, ok, "assistant message is not a map")
	toolCalls, ok := assistant["tool_calls"].([]any)
	require.True(t, ok)
	require.Len(t, toolCalls, 1)

	call, ok := toolCalls[0].(map[string]any)
	require.True(t, ok, "tool_call entry is not a map")
	assert.Equal(t, "call_9", call["id"])
	assert.Equal(t, "function", call["type"])

	toolMsg, ok := msgs[2].(map[string]any)
	require.True(t, ok, "tool message is not a map")
	assert.Equal(t, "call_9", toolMsg["tool_call_id"])
}

func TestAzureClient_URLConstruction_TrimsTrailingSlashFromEndpoint(t *testing.T) {
	var capturedPath string

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedPath = r.URL.Path
		writeJSON(t, w, validAzureResponse("ok"))
	}))
	defer server.Close()

	client := ai.NewAzureClient(server.URL+"/", "test-key", "gpt-4o-mini", "2024-12-01-preview")
	_, err := client.Generate(context.Background(), providers.GenerateParams{
		SystemPrompt: "s",
		UserPrompt:   "u",
	})
	require.NoError(t, err)

	assert.Equal(t, "/openai/deployments/gpt-4o-mini/chat/completions", capturedPath)
}
