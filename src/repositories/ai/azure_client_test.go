package ai_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

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

// TestAzureClient_Generate cubre los tres casos del método Generate.
func TestAzureClient_Generate(t *testing.T) {
	t.Run("returns generated content", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Header.Get("api-key") != "test-key" {
				t.Errorf("expected api-key header to be 'test-key', got %q", r.Header.Get("api-key"))
			}
			if r.Header.Get("Content-Type") != "application/json" {
				t.Errorf("expected Content-Type application/json, got %q", r.Header.Get("Content-Type"))
			}
			writeJSON(t, w, validAzureResponse("test response"))
		}))
		defer server.Close()

		client := ai.NewAzureClient(server.URL, "test-key", "gpt-4o-mini")
		got, err := client.Generate(context.Background(), providers.GenerateParams{
			SystemPrompt: "you are helpful",
			UserPrompt:   "hello",
		})

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if got != "test response" {
			t.Errorf("expected 'test response', got %q", got)
		}
	})

	t.Run("returns error for non-200 status", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			http.Error(w, "internal server error", http.StatusInternalServerError)
		}))
		defer server.Close()

		client := ai.NewAzureClient(server.URL, "test-key", "gpt-4o-mini")
		_, err := client.Generate(context.Background(), providers.GenerateParams{
			SystemPrompt: "you are helpful",
			UserPrompt:   "hello",
		})

		if err == nil {
			t.Fatal("expected error for 500 status, got nil")
		}
		if !strings.Contains(err.Error(), "500") {
			t.Errorf("expected error to mention status 500, got: %v", err)
		}
	})

	t.Run("returns error for empty choices", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			writeJSON(t, w, map[string]any{"choices": []any{}})
		}))
		defer server.Close()

		client := ai.NewAzureClient(server.URL, "test-key", "gpt-4o-mini")
		_, err := client.Generate(context.Background(), providers.GenerateParams{
			SystemPrompt: "you are helpful",
			UserPrompt:   "hello",
		})

		if err == nil {
			t.Fatal("expected error for empty choices, got nil")
		}
		if !strings.Contains(err.Error(), "empty response") {
			t.Errorf("expected error to mention 'empty response', got: %v", err)
		}
	})
}

// TestAzureClient_Chat cubre los casos del método Chat.
func TestAzureClient_Chat(t *testing.T) {
	t.Run("returns chat response", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			writeJSON(t, w, validAzureResponse("chat reply"))
		}))
		defer server.Close()

		client := ai.NewAzureClient(server.URL, "test-key", "gpt-4o-mini")
		resp, err := client.Chat(context.Background(), []providers.ChatMessage{
			{Role: "user", Content: "hello"},
		})

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if resp == nil {
			t.Fatal("expected non-nil response")
		}
		if resp.Content != "chat reply" {
			t.Errorf("expected 'chat reply', got %q", resp.Content)
		}
	})

	t.Run("sends messages correctly", func(t *testing.T) {
		var capturedBody map[string]any

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if err := json.NewDecoder(r.Body).Decode(&capturedBody); err != nil {
				t.Errorf("decode request body: %v", err)
			}
			writeJSON(t, w, validAzureResponse("ok"))
		}))
		defer server.Close()

		client := ai.NewAzureClient(server.URL, "test-key", "gpt-4o-mini")
		messages := []providers.ChatMessage{
			{Role: "system", Content: "be helpful"},
			{Role: "user", Content: "what is Go?"},
		}
		if _, err := client.Chat(context.Background(), messages); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		rawMessages, ok := capturedBody["messages"].([]any)
		if !ok {
			t.Fatal("messages field missing or not a slice in request body")
		}
		if len(rawMessages) != 2 {
			t.Fatalf("expected 2 messages, got %d", len(rawMessages))
		}

		first, ok := rawMessages[0].(map[string]any)
		if !ok {
			t.Fatal("first message is not a map")
		}
		if first["role"] != "system" {
			t.Errorf("expected role 'system', got %v", first["role"])
		}
		if first["content"] != "be helpful" {
			t.Errorf("expected content 'be helpful', got %v", first["content"])
		}

		second, ok := rawMessages[1].(map[string]any)
		if !ok {
			t.Fatal("second message is not a map")
		}
		if second["role"] != "user" {
			t.Errorf("expected role 'user', got %v", second["role"])
		}
		if second["content"] != "what is Go?" {
			t.Errorf("expected content 'what is Go?', got %v", second["content"])
		}
	})
}

// TestAzureClient_ChatWithTools verifica que ChatWithTools delega a doRequest correctamente.
func TestAzureClient_ChatWithTools(t *testing.T) {
	t.Run("returns response", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			writeJSON(t, w, validAzureResponse("tool reply"))
		}))
		defer server.Close()

		client := ai.NewAzureClient(server.URL, "test-key", "gpt-4o-mini")
		resp, err := client.ChatWithTools(
			context.Background(),
			[]providers.ChatMessage{{Role: "user", Content: "use a tool"}},
			[]providers.ToolDefinition{{Name: "my_tool", Description: "does stuff"}},
		)

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if resp == nil {
			t.Fatal("expected non-nil response")
		}
		if resp.Content != "tool reply" {
			t.Errorf("expected 'tool reply', got %q", resp.Content)
		}
	})

	t.Run("sends tool definitions in openai function format", func(t *testing.T) {
		var capturedBody map[string]any
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if err := json.NewDecoder(r.Body).Decode(&capturedBody); err != nil {
				t.Errorf("decode request body: %v", err)
			}
			writeJSON(t, w, validAzureResponse("ok"))
		}))
		defer server.Close()

		client := ai.NewAzureClient(server.URL, "test-key", "gpt-4o-mini")
		_, err := client.ChatWithTools(
			context.Background(),
			[]providers.ChatMessage{{Role: "user", Content: "buscar dispositivos"}},
			[]providers.ToolDefinition{{
				Name:        "search_devices",
				Description: "busca dispositivos",
				Parameters:  map[string]any{"type": "object"},
			}},
		)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		tools, ok := capturedBody["tools"].([]any)
		if !ok || len(tools) != 1 {
			t.Fatalf("expected 1 tool in request body, got %v", capturedBody["tools"])
		}
		tool, ok := tools[0].(map[string]any)
		if !ok {
			t.Fatal("tool is not a map")
		}
		if tool["type"] != "function" {
			t.Errorf("expected tool type 'function', got %v", tool["type"])
		}
		fn, ok := tool["function"].(map[string]any)
		if !ok {
			t.Fatal("tool.function is not a map")
		}
		if fn["name"] != "search_devices" {
			t.Errorf("expected function name 'search_devices', got %v", fn["name"])
		}
	})

	t.Run("parses tool_calls from response", func(t *testing.T) {
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

		client := ai.NewAzureClient(server.URL, "test-key", "gpt-4o-mini")
		resp, err := client.ChatWithTools(
			context.Background(),
			[]providers.ChatMessage{{Role: "user", Content: "buscar timer"}},
			[]providers.ToolDefinition{{Name: "search_devices", Description: "busca"}},
		)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(resp.ToolCalls) != 1 {
			t.Fatalf("expected 1 tool call, got %d", len(resp.ToolCalls))
		}
		tc := resp.ToolCalls[0]
		if tc.ID != "call_1" || tc.Name != "search_devices" {
			t.Errorf("unexpected tool call: %+v", tc)
		}
		if tc.Arguments != `{"query":"timer"}` {
			t.Errorf("unexpected arguments: %q", tc.Arguments)
		}
	})

	t.Run("serializes assistant tool_calls and tool result messages in openai format", func(t *testing.T) {
		var capturedBody map[string]any
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if err := json.NewDecoder(r.Body).Decode(&capturedBody); err != nil {
				t.Errorf("decode request body: %v", err)
			}
			writeJSON(t, w, validAzureResponse("ok"))
		}))
		defer server.Close()

		client := ai.NewAzureClient(server.URL, "test-key", "gpt-4o-mini")
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
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		msgs, ok := capturedBody["messages"].([]any)
		if !ok || len(msgs) != 3 {
			t.Fatalf("expected 3 messages, got %v", capturedBody["messages"])
		}

		assistant, ok := msgs[1].(map[string]any)
		if !ok {
			t.Fatal("assistant message is not a map")
		}
		toolCalls, ok := assistant["tool_calls"].([]any)
		if !ok || len(toolCalls) != 1 {
			t.Fatalf("expected assistant.tool_calls to have 1 entry, got %v", assistant["tool_calls"])
		}
		call, ok := toolCalls[0].(map[string]any)
		if !ok {
			t.Fatal("tool_call entry is not a map")
		}
		if call["id"] != "call_9" || call["type"] != "function" {
			t.Errorf("unexpected tool_call envelope: %v", call)
		}

		toolMsg, ok := msgs[2].(map[string]any)
		if !ok {
			t.Fatal("tool message is not a map")
		}
		if toolMsg["tool_call_id"] != "call_9" {
			t.Errorf("expected tool_call_id 'call_9', got %v", toolMsg["tool_call_id"])
		}
	})
}

// TestAzureClient_URLConstruction verifica que el endpoint se construye correctamente
// incluyendo el trimado de slashes finales.
func TestAzureClient_URLConstruction(t *testing.T) {
	t.Run("trims trailing slash from endpoint", func(t *testing.T) {
		var capturedPath string

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			capturedPath = r.URL.Path
			writeJSON(t, w, validAzureResponse("ok"))
		}))
		defer server.Close()

		// Endpoint con slash final — debe ser recortado para no generar doble slash.
		client := ai.NewAzureClient(server.URL+"/", "test-key", "gpt-4o-mini")
		if _, err := client.Generate(context.Background(), providers.GenerateParams{
			SystemPrompt: "s",
			UserPrompt:   "u",
		}); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		expected := "/openai/deployments/gpt-4o-mini/chat/completions"
		if capturedPath != expected {
			t.Errorf("expected path %q, got %q", expected, capturedPath)
		}
	})
}
