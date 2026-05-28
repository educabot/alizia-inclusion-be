package inclusion

import (
	"strings"
	"testing"

	"github.com/educabot/alizia-inclusion-be/src/core/providers"
)

func TestCapMessages_ReturnsMessagesUnchangedWhenUnderBudget(t *testing.T) {
	messages := []providers.ChatMessage{
		{Role: "system", Content: "system prompt"},
		{Role: "user", Content: "hola"},
		{Role: "assistant", Content: "buenas"},
		{Role: "user", Content: "una pregunta"},
	}

	got := capMessages(messages, defaultMaxHistoryTokens)

	if len(got) != len(messages) {
		t.Fatalf("expected %d messages, got %d", len(messages), len(got))
	}
}

func TestCapMessages_DropsOldestHistoryButKeepsSystemAndCurrentTurnWhenOverBudget(t *testing.T) {
	big := strings.Repeat("x", 4000) // ~1000 tokens each
	messages := []providers.ChatMessage{
		{Role: "system", Content: "system"},
		{Role: "user", Content: "OLDEST " + big},
		{Role: "assistant", Content: "MIDDLE " + big},
		{Role: "user", Content: "RECENT " + big},
		{Role: "user", Content: "CURRENT question"},
	}

	got := capMessages(messages, 2500)

	if got[0].Role != "system" {
		t.Errorf("expected first message to be system, got %q", got[0].Role)
	}
	if got[len(got)-1].Content != "CURRENT question" {
		t.Errorf("expected last message to be the current turn, got %q", got[len(got)-1].Content)
	}
	for _, m := range got {
		if strings.HasPrefix(m.Content, "OLDEST") {
			t.Error("expected oldest history message to be dropped")
		}
	}
}

func TestCapMessages_KeepsOnlySystemAndCurrentTurnWhenSystemPromptExhaustsBudget(t *testing.T) {
	huge := strings.Repeat("y", 12000) // ~3000 tokens
	messages := []providers.ChatMessage{
		{Role: "system", Content: huge},
		{Role: "user", Content: strings.Repeat("z", 4000)},
		{Role: "user", Content: "current"},
	}

	got := capMessages(messages, 2000)

	if len(got) != 2 {
		t.Fatalf("expected 2 messages (system + current), got %d", len(got))
	}
	if got[1].Content != "current" {
		t.Errorf("expected current turn preserved, got %q", got[1].Content)
	}
}

func TestCapMessages_ReturnsShortConversationsUntouched(t *testing.T) {
	messages := []providers.ChatMessage{
		{Role: "system", Content: strings.Repeat("a", 100000)},
		{Role: "user", Content: strings.Repeat("b", 100000)},
	}

	got := capMessages(messages, 10)

	if len(got) != 2 {
		t.Errorf("expected 2 messages preserved, got %d", len(got))
	}
}

func TestEstimateTokens(t *testing.T) {
	tests := []struct {
		input string
		want  int
	}{
		{"", 0},
		{"abcd", 1},
		{"abcde", 2},
	}
	for _, tt := range tests {
		got := estimateTokens(tt.input)
		if got != tt.want {
			t.Errorf("estimateTokens(%q) = %d, want %d", tt.input, got, tt.want)
		}
	}
}
