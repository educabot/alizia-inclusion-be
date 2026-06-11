package inclusion

import (
	"context"
	"log/slog"

	"github.com/google/uuid"

	"github.com/educabot/alizia-inclusion-be/src/core/providers"
)

// aiTrace captures one AI turn for usage recording and trace (HU-6, T-6.5).
// ContextSnapshot must contain only IDs (no PII): student_id, classroom_id, etc.
type aiTrace struct {
	orgID          uuid.UUID
	userID         int64
	mode           string
	model          string
	latencyMs      int
	toolCalls      int
	conversationID int64
	context        map[string]any
	usage          *providers.TokenUsage
}

// recordAIUsage persists the usage record and trace for one AI turn (HU-6, T-6.5).
// Every turn is traced even when the provider reports no token counts — model, latency,
// tool_calls, and context are still valuable. Best-effort: a nil provider or anonymous
// request is skipped, and a record failure is logged rather than propagated so it never
// blocks the teacher's response.
func recordAIUsage(ctx context.Context, usage providers.AIUsageProvider, t aiTrace) {
	if usage == nil || t.userID == 0 {
		return
	}
	record := providers.AIUsageRecord{
		OrgID:           t.orgID,
		UserID:          t.userID,
		Mode:            t.mode,
		Model:           t.model,
		LatencyMs:       t.latencyMs,
		ToolCalls:       t.toolCalls,
		ContextSnapshot: t.context,
	}
	if t.usage != nil {
		record.PromptTokens = t.usage.PromptTokens
		record.CompletionTokens = t.usage.CompletionTokens
		record.TotalTokens = t.usage.TotalTokens
	}
	if t.conversationID > 0 {
		record.ConversationID = &t.conversationID
	}
	if err := usage.Record(ctx, record); err != nil {
		slog.WarnContext(ctx, "record ai usage failed", "error", err, "user_id", t.userID, "mode", t.mode)
	}
}
