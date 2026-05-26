package dashboard

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/educabot/alizia-inclusion-be/src/core/providers"
)

// defaultAIUsageWindowDays is the look-back window used when the caller does not
// specify one. Expressed in days.
const defaultAIUsageWindowDays = 30

// maxAIUsageWindowDays caps how far back a single query may look, to bound the
// aggregation cost. Expressed in days.
const maxAIUsageWindowDays = 365

type GetAIUsageRequest struct {
	OrgID uuid.UUID
	// Days is the look-back window. Zero or negative falls back to the default.
	Days int
}

func (r GetAIUsageRequest) Validate() error {
	if r.OrgID == uuid.Nil {
		return errOrgIDRequired
	}
	return nil
}

type AIUsageModeResponse struct {
	Mode             string `json:"mode"`
	Requests         int    `json:"requests"`
	PromptTokens     int    `json:"prompt_tokens"`
	CompletionTokens int    `json:"completion_tokens"`
	TotalTokens      int    `json:"total_tokens"`
}

type GetAIUsageResponse struct {
	WindowDays       int                   `json:"window_days"`
	TotalRequests    int                   `json:"total_requests"`
	PromptTokens     int                   `json:"prompt_tokens"`
	CompletionTokens int                   `json:"completion_tokens"`
	TotalTokens      int                   `json:"total_tokens"`
	ByMode           []AIUsageModeResponse `json:"by_mode"`
}

type GetAIUsage interface {
	Execute(ctx context.Context, req GetAIUsageRequest) (*GetAIUsageResponse, error)
}

type getAIUsageImpl struct {
	usage providers.AIUsageProvider
}

func NewGetAIUsage(usage providers.AIUsageProvider) GetAIUsage {
	return &getAIUsageImpl{usage: usage}
}

func (uc *getAIUsageImpl) Execute(ctx context.Context, req GetAIUsageRequest) (*GetAIUsageResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	days := req.Days
	if days <= 0 {
		days = defaultAIUsageWindowDays
	}
	if days > maxAIUsageWindowDays {
		days = maxAIUsageWindowDays
	}

	since := time.Now().AddDate(0, 0, -days)
	summary, err := uc.usage.Summarize(ctx, req.OrgID, since)
	if err != nil {
		return nil, err
	}

	byMode := make([]AIUsageModeResponse, len(summary.ByMode))
	for i, m := range summary.ByMode {
		byMode[i] = AIUsageModeResponse{
			Mode:             m.Mode,
			Requests:         m.Requests,
			PromptTokens:     m.PromptTokens,
			CompletionTokens: m.CompletionTokens,
			TotalTokens:      m.TotalTokens,
		}
	}

	return &GetAIUsageResponse{
		WindowDays:       days,
		TotalRequests:    summary.TotalRequests,
		PromptTokens:     summary.PromptTokens,
		CompletionTokens: summary.CompletionTokens,
		TotalTokens:      summary.TotalTokens,
		ByMode:           byMode,
	}, nil
}
