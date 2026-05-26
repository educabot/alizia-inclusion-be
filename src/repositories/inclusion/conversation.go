package inclusion

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
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

func (r *conversationRepo) AppendTurn(ctx context.Context, params providers.AppendTurnParams) (int64, error) {
	var convID int64
	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if params.ConversationID == 0 {
			conv := entities.Conversation{
				OrganizationID: params.OrgID,
				UserID:         params.UserID,
				StudentID:      params.StudentID,
				Mode:           params.Mode,
			}
			if err := tx.Create(&conv).Error; err != nil {
				return fmt.Errorf("create conversation: %w", err)
			}
			convID = conv.ID
		} else {
			res := tx.Model(&entities.Conversation{}).
				Where("id = ? AND organization_id = ? AND user_id = ?", params.ConversationID, params.OrgID, params.UserID).
				Update("updated_at", time.Now())
			if res.Error != nil {
				return fmt.Errorf("touch conversation: %w", res.Error)
			}
			if res.RowsAffected == 0 {
				return providers.ErrNotFound
			}
			convID = params.ConversationID
		}

		metaJSON := []byte("{}")
		if len(params.Metadata) > 0 {
			b, err := json.Marshal(params.Metadata)
			if err != nil {
				return fmt.Errorf("marshal metadata: %w", err)
			}
			metaJSON = b
		}

		msgs := []entities.ConversationMessage{
			{
				ConversationID: convID,
				Role:           "user",
				Content:        params.UserContent,
				Metadata:       datatypes.JSON([]byte("{}")),
			},
			{
				ConversationID: convID,
				Role:           "assistant",
				Content:        params.AssistantContent,
				Metadata:       datatypes.JSON(metaJSON),
			},
		}
		if err := tx.Create(&msgs).Error; err != nil {
			return fmt.Errorf("create messages: %w", err)
		}
		return nil
	})
	if err != nil {
		return 0, err
	}
	return convID, nil
}
