package inclusion

import (
	"context"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/educabot/alizia-inclusion-be/src/core/entities"
	"github.com/educabot/alizia-inclusion-be/src/core/providers"
)

type aiUsageRepo struct {
	db *gorm.DB
}

func NewAIUsageRepo(db *gorm.DB) providers.AIUsageProvider {
	return &aiUsageRepo{db: db}
}

func (r *aiUsageRepo) Record(ctx context.Context, record providers.AIUsageRecord) error {
	usage := entities.AIUsage{
		OrganizationID:   record.OrgID,
		UserID:           record.UserID,
		Mode:             record.Mode,
		PromptTokens:     record.PromptTokens,
		CompletionTokens: record.CompletionTokens,
		TotalTokens:      record.TotalTokens,
	}
	return r.db.WithContext(ctx).Create(&usage).Error
}

func (r *aiUsageRepo) Summarize(ctx context.Context, orgID uuid.UUID, since time.Time) (*providers.AIUsageSummary, error) {
	var rows []struct {
		Mode             string
		Requests         int
		PromptTokens     int
		CompletionTokens int
		TotalTokens      int
	}
	err := r.db.WithContext(ctx).
		Model(&entities.AIUsage{}).
		Select("mode, "+
			"COUNT(*) AS requests, "+
			"COALESCE(SUM(prompt_tokens), 0) AS prompt_tokens, "+
			"COALESCE(SUM(completion_tokens), 0) AS completion_tokens, "+
			"COALESCE(SUM(total_tokens), 0) AS total_tokens").
		Where("organization_id = ? AND created_at >= ?", orgID, since).
		Group("mode").
		Order("total_tokens DESC").
		Scan(&rows).Error
	if err != nil {
		return nil, err
	}

	summary := &providers.AIUsageSummary{ByMode: make([]providers.AIUsageModeSummary, 0, len(rows))}
	for _, row := range rows {
		summary.TotalRequests += row.Requests
		summary.PromptTokens += row.PromptTokens
		summary.CompletionTokens += row.CompletionTokens
		summary.TotalTokens += row.TotalTokens
		summary.ByMode = append(summary.ByMode, providers.AIUsageModeSummary{
			Mode:             row.Mode,
			Requests:         row.Requests,
			PromptTokens:     row.PromptTokens,
			CompletionTokens: row.CompletionTokens,
			TotalTokens:      row.TotalTokens,
		})
	}
	return summary, nil
}
