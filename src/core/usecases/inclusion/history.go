package inclusion

import "github.com/educabot/alizia-inclusion-be/src/core/providers"

// defaultMaxHistoryTokens bounds the estimated token budget for an AI request.
// It leaves headroom under typical context windows for the model's completion.
const defaultMaxHistoryTokens = 3000

// estimateTokens approximates the token count of a string using the common
// heuristic of ~4 characters per token. It is intentionally cheap and rough —
// good enough to decide what to drop before sending to the model.
func estimateTokens(s string) int {
	return (len(s) + 3) / 4
}

// capMessages trims a conversation so its estimated token total stays under
// maxTokens. The system prompt (first message) and the current turn (last
// message) are always preserved; the oldest history messages in between are
// dropped first, keeping the most recent context.
func capMessages(messages []providers.ChatMessage, maxTokens int) []providers.ChatMessage {
	if len(messages) <= 2 {
		return messages
	}

	total := 0
	for i := range messages {
		total += estimateTokens(messages[i].Content)
	}
	if total <= maxTokens {
		return messages
	}

	system := messages[0]
	last := messages[len(messages)-1]
	middle := messages[1 : len(messages)-1]

	budget := maxTokens - estimateTokens(system.Content) - estimateTokens(last.Content)

	kept := make([]providers.ChatMessage, 0, len(middle))
	used := 0
	for i := len(middle) - 1; i >= 0; i-- {
		t := estimateTokens(middle[i].Content)
		if used+t > budget {
			break
		}
		used += t
		kept = append([]providers.ChatMessage{middle[i]}, kept...)
	}

	out := make([]providers.ChatMessage, 0, len(kept)+2)
	out = append(out, system)
	out = append(out, kept...)
	out = append(out, last)
	return out
}
