package providers

import (
	"context"

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

type AIUsageProvider interface {
	Record(ctx context.Context, record AIUsageRecord) error
}
