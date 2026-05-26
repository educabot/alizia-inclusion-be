package inclusion

import (
	"context"

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
