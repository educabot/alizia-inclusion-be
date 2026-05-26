package inclusion

import (
	"context"
	"log/slog"

	"github.com/google/uuid"

	"github.com/educabot/alizia-inclusion-be/src/core/providers"
)

// recordAIUsage persists token usage for an AI call. It is best-effort: a nil
// provider, missing usage data, or an anonymous request are skipped, and a
// recording failure is logged rather than propagated so it never blocks the
// user-facing response.
func recordAIUsage(ctx context.Context, usage providers.AIUsageProvider, orgID uuid.UUID, userID int64, mode string, tokens *providers.TokenUsage) {
	if usage == nil || tokens == nil || userID == 0 {
		return
	}
	err := usage.Record(ctx, providers.AIUsageRecord{
		OrgID:            orgID,
		UserID:           userID,
		Mode:             mode,
		PromptTokens:     tokens.PromptTokens,
		CompletionTokens: tokens.CompletionTokens,
		TotalTokens:      tokens.TotalTokens,
	})
	if err != nil {
		slog.WarnContext(ctx, "record ai usage failed", "error", err, "user_id", userID, "mode", mode)
	}
}
