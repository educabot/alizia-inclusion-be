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

// ListPendingSummary trae conversaciones "cerradas" (último mensaje anterior a
// idleBefore) cuyo resumen falta o quedó viejo respecto del último mensaje, con
// sus Messages precargados en orden. Lo consume el batch de resúmenes (cron).
func (r *conversationRepo) ListPendingSummary(ctx context.Context, idleBefore time.Time, limit int) ([]entities.Conversation, error) {
	var ids []int64
	err := r.db.WithContext(ctx).Raw(`
		SELECT c.id
		FROM conversations c
		JOIN conversation_messages m ON m.conversation_id = c.id
		LEFT JOIN conversation_summaries s ON s.conversation_id = c.id
		GROUP BY c.id, s.updated_at
		HAVING MAX(m.created_at) < ?
		   AND (s.updated_at IS NULL OR s.updated_at < MAX(m.created_at))
		ORDER BY MAX(m.created_at) ASC
		LIMIT ?`, idleBefore, limit).Scan(&ids).Error
	if err != nil {
		return nil, err
	}
	if len(ids) == 0 {
		return nil, nil
	}

	var conversations []entities.Conversation
	err = r.db.WithContext(ctx).
		Preload("Messages", func(db *gorm.DB) *gorm.DB {
			return db.Order("created_at ASC")
		}).
		Where("id IN ?", ids).
		Find(&conversations).Error
	if err != nil {
		return nil, err
	}
	return conversations, nil
}

func (r *conversationRepo) Delete(ctx context.Context, orgID uuid.UUID, id int64) error {
	result := r.db.WithContext(ctx).
		Where("organization_id = ? AND id = ?", orgID, id).
		Delete(&entities.Conversation{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return providers.ErrNotFound
	}
	return nil
}

func (r *conversationRepo) Rename(ctx context.Context, orgID uuid.UUID, id int64, title string) error {
	result := r.db.WithContext(ctx).
		Model(&entities.Conversation{}).
		Where("organization_id = ? AND id = ?", orgID, id).
		Update("title", title)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return providers.ErrNotFound
	}
	return nil
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
