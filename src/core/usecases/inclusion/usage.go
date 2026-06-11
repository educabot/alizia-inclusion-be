package inclusion

import (
	"context"
	"log/slog"

	"github.com/google/uuid"

	"github.com/educabot/alizia-inclusion-be/src/core/providers"
)

// aiTrace describe un turno de IA para el registro de uso + traza (HU-6, T-6.5).
// ContextSnapshot debe contener solo IDs (sin PII): student_id, classroom_id, etc.
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

// recordAIUsage persiste el uso + la traza de un turno de IA. Es best-effort: un
// provider nil, sin datos de tokens o un request anónimo se saltean, y un fallo
// de registro se loguea en vez de propagarse, para que nunca bloquee la respuesta
// al docente.
func recordAIUsage(ctx context.Context, usage providers.AIUsageProvider, t aiTrace) {
	if usage == nil || t.usage == nil || t.userID == 0 {
		return
	}
	record := providers.AIUsageRecord{
		OrgID:            t.orgID,
		UserID:           t.userID,
		Mode:             t.mode,
		PromptTokens:     t.usage.PromptTokens,
		CompletionTokens: t.usage.CompletionTokens,
		TotalTokens:      t.usage.TotalTokens,
		Model:            t.model,
		LatencyMs:        t.latencyMs,
		ToolCalls:        t.toolCalls,
		ContextSnapshot:  t.context,
	}
	if t.conversationID > 0 {
		record.ConversationID = &t.conversationID
	}
	if err := usage.Record(ctx, record); err != nil {
		slog.WarnContext(ctx, "record ai usage failed", "error", err, "user_id", t.userID, "mode", t.mode)
	}
}
