package inclusion

import (
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/educabot/alizia-inclusion-be/src/core/entities"
	"github.com/educabot/alizia-inclusion-be/src/core/providers"
)

type conversationSummaryRepo struct {
	db *gorm.DB
}

func NewConversationSummaryRepo(db *gorm.DB) providers.ConversationSummaryProvider {
	return &conversationSummaryRepo{db: db}
}

// RecentByStudent trae los resúmenes más recientes ligados a un alumno, acotados por org.
func (r *conversationSummaryRepo) RecentByStudent(ctx context.Context, orgID uuid.UUID, studentID int64, limit int) ([]entities.ConversationSummary, error) {
	var out []entities.ConversationSummary
	err := r.db.WithContext(ctx).
		Model(&entities.ConversationSummary{}).
		Select("conversation_summaries.*").
		Joins("JOIN conversation_summary_students css ON css.conversation_id = conversation_summaries.conversation_id").
		Joins("JOIN conversations c ON c.id = conversation_summaries.conversation_id").
		Where("css.student_id = ? AND c.organization_id = ?", studentID, orgID).
		Order("conversation_summaries.updated_at DESC").
		Limit(limit).
		Find(&out).Error
	if err != nil {
		return nil, err
	}
	return out, nil
}

// RecentByDevice trae los resúmenes más recientes ligados a un dispositivo de la valija.
func (r *conversationSummaryRepo) RecentByDevice(ctx context.Context, orgID uuid.UUID, deviceID int64, limit int) ([]entities.ConversationSummary, error) {
	var out []entities.ConversationSummary
	err := r.db.WithContext(ctx).
		Model(&entities.ConversationSummary{}).
		Select("conversation_summaries.*").
		Joins("JOIN conversation_summary_devices csd ON csd.conversation_id = conversation_summaries.conversation_id").
		Joins("JOIN conversations c ON c.id = conversation_summaries.conversation_id").
		Where("csd.device_id = ? AND c.organization_id = ?", deviceID, orgID).
		Order("conversation_summaries.updated_at DESC").
		Limit(limit).
		Find(&out).Error
	if err != nil {
		return nil, err
	}
	return out, nil
}

// RecentByTopic trae los resúmenes más recientes cuyo topic_keywords contiene la keyword.
func (r *conversationSummaryRepo) RecentByTopic(ctx context.Context, orgID uuid.UUID, keyword string, limit int) ([]entities.ConversationSummary, error) {
	var out []entities.ConversationSummary
	err := r.db.WithContext(ctx).
		Model(&entities.ConversationSummary{}).
		Select("conversation_summaries.*").
		Joins("JOIN conversations c ON c.id = conversation_summaries.conversation_id").
		Where("c.organization_id = ? AND ? = ANY(conversation_summaries.topic_keywords)", orgID, keyword).
		Order("conversation_summaries.updated_at DESC").
		Limit(limit).
		Find(&out).Error
	if err != nil {
		return nil, err
	}
	return out, nil
}
