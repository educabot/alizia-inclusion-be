package ai_test

import (
	"context"
	"strings"
	"testing"

	ai "github.com/educabot/alizia-inclusion-be/src/repositories/ai"

	"github.com/educabot/alizia-inclusion-be/src/core/providers"
)

func TestStubClient_Generate(t *testing.T) {
	t.Run("returns stub content", func(t *testing.T) {
		client := ai.NewStubClient()
		got, err := client.Generate(context.Background(), providers.GenerateParams{
			UserPrompt: "test prompt",
		})

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if got == "" {
			t.Error("expected non-empty content")
		}
		if !strings.Contains(got, "[stub]") {
			t.Errorf("expected content to contain '[stub]', got %q", got)
		}
	})
}

func TestStubClient_Chat(t *testing.T) {
	t.Run("returns stub response", func(t *testing.T) {
		client := ai.NewStubClient()
		resp, err := client.Chat(context.Background(), []providers.ChatMessage{
			{Role: "user", Content: "hello"},
		})

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if resp == nil {
			t.Fatal("expected non-nil response")
		}
		if resp.Content == "" {
			t.Error("expected non-empty response content")
		}
	})
}

func TestStubClient_ChatWithTools(t *testing.T) {
	t.Run("returns stub response", func(t *testing.T) {
		client := ai.NewStubClient()
		resp, err := client.ChatWithTools(
			context.Background(),
			[]providers.ChatMessage{{Role: "user", Content: "use tool"}},
			[]providers.ToolDefinition{{Name: "noop", Description: "does nothing"}},
		)

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if resp == nil {
			t.Fatal("expected non-nil response")
		}
	})
}
