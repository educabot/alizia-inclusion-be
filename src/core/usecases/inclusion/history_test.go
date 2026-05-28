package inclusion

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

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

	assert.Len(t, got, len(messages))
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

	assert.Equal(t, "system", got[0].Role)
	assert.Equal(t, "CURRENT question", got[len(got)-1].Content)
	for _, m := range got {
		assert.False(t, strings.HasPrefix(m.Content, "OLDEST"), "oldest history message should be dropped")
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

	assert.Len(t, got, 2)
	assert.Equal(t, "current", got[1].Content)
}

func TestCapMessages_ReturnsShortConversationsUntouched(t *testing.T) {
	messages := []providers.ChatMessage{
		{Role: "system", Content: strings.Repeat("a", 100000)},
		{Role: "user", Content: strings.Repeat("b", 100000)},
	}

	got := capMessages(messages, 10)

	assert.Len(t, got, 2)
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
		assert.Equal(t, tt.want, got)
	}
}
