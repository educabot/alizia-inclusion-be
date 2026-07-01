package inclusion

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/educabot/alizia-inclusion-be/src/core/entities"
	"github.com/educabot/alizia-inclusion-be/src/core/providers"
)

type messageFeedbackRepo struct {
	db *gorm.DB
}

func NewMessageFeedbackRepo(db *gorm.DB) providers.MessageFeedbackProvider {
	return &messageFeedbackRepo{db: db}
}

func (r *messageFeedbackRepo) MessageContext(ctx context.Context, messageID int64) (int64, uuid.UUID, error) {
	var row struct {
		ConversationID int64
		OrganizationID uuid.UUID
	}
	err := r.db.WithContext(ctx).
		Table("conversation_messages AS cm").
		Select("c.id AS conversation_id, c.organization_id AS organization_id").
		Joins("JOIN conversations c ON c.id = cm.conversation_id").
		Where("cm.id = ?", messageID).
		Take(&row).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return 0, uuid.Nil, providers.ErrNotFound
	}
	if err != nil {
		return 0, uuid.Nil, err
	}
	return row.ConversationID, row.OrganizationID, nil
}

func (r *messageFeedbackRepo) Upsert(ctx context.Context, fb *entities.MessageFeedback) error {
	now := time.Now()
	fb.UpdatedAt = now
	return r.db.WithContext(ctx).
		Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "conversation_message_id"}, {Name: "user_id"}},
			DoUpdates: clause.Assignments(map[string]any{"rating": fb.Rating, "comment": fb.Comment, "updated_at": now}),
		}).
		Create(fb).Error
}

func (r *messageFeedbackRepo) Delete(ctx context.Context, orgID uuid.UUID, messageID, userID int64) error {
	return r.db.WithContext(ctx).
		Where("organization_id = ? AND conversation_message_id = ? AND user_id = ?", orgID, messageID, userID).
		Delete(&entities.MessageFeedback{}).Error
}

func (r *messageFeedbackRepo) List(ctx context.Context, orgID uuid.UUID, rating string) ([]providers.MessageFeedbackReview, error) {
	// Traemos el feedback + el contenido del mensaje comentado + la pregunta del
	// usuario inmediatamente anterior (último mensaje 'user' con id menor en la misma
	// conversación), para revisar el error con contexto.
	query := r.db.WithContext(ctx).
		Table("message_feedback AS mf").
		Select(`mf.*, cm.content AS message_content,
			(SELECT pm.content FROM conversation_messages pm
				WHERE pm.conversation_id = mf.conversation_id
					AND pm.id < mf.conversation_message_id
					AND pm.role = 'user'
				ORDER BY pm.id DESC LIMIT 1) AS previous_user_message`).
		Joins("JOIN conversation_messages cm ON cm.id = mf.conversation_message_id").
		Where("mf.organization_id = ?", orgID)
	if rating != "" {
		query = query.Where("mf.rating = ?", rating)
	}

	var out []providers.MessageFeedbackReview
	if err := query.Order("mf.created_at DESC").Scan(&out).Error; err != nil {
		return nil, err
	}
	return out, nil
}
