package providers

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type AIUsageRecord struct {
	OrgID            uuid.UUID
	UserID           int64
	Mode             string
	PromptTokens     int
	CompletionTokens int
	TotalTokens      int
}

// AIUsageModeSummary aggregates token consumption for a single mode (assist,
// recommend, etc.) over a time window.
type AIUsageModeSummary struct {
	Mode             string
	Requests         int
	PromptTokens     int
	CompletionTokens int
	TotalTokens      int
}

// AIUsageSummary aggregates token consumption for an organization over a window.
type AIUsageSummary struct {
	TotalRequests    int
	PromptTokens     int
	CompletionTokens int
	TotalTokens      int
	ByMode           []AIUsageModeSummary
}

type AIUsageProvider interface {
	Record(ctx context.Context, record AIUsageRecord) error
	// Summarize aggregates usage for an org since the given time, grouped by mode.
	Summarize(ctx context.Context, orgID uuid.UUID, since time.Time) (*AIUsageSummary, error)
}
