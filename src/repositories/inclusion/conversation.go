package inclusion

import (
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/educabot/alizia-inclusion-be/src/core/entities"
	"github.com/educabot/alizia-inclusion-be/src/core/providers"
)

type conversationRepo struct {
	db *gorm.DB
}

func NewConversationRepo(db *gorm.DB) providers.ConversationProvider {
	return &conversationRepo{db: db}
}

func (r *conversationRepo) ListByUser(ctx context.Context, orgID uuid.UUID, userID int64, mode string) ([]entities.Conversation, error) {
	var conversations []entities.Conversation
	q := r.db.WithContext(ctx).
		Preload("Messages", func(db *gorm.DB) *gorm.DB {
			return db.Order("created_at ASC")
		}).
		Where("organization_id = ? AND user_id = ?", orgID, userID)
	if mode != "" {
		q = q.Where("mode = ?", mode)
	}
	err := q.Order("updated_at DESC").Find(&conversations).Error
	if err != nil {
		return nil, err
	}
	return conversations, nil
}
