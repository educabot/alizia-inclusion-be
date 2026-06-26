package inclusion

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/educabot/alizia-inclusion-be/src/core/entities"
	"github.com/educabot/alizia-inclusion-be/src/core/providers"
)

type conversationSummaryRepo struct {
	db *gorm.DB
}

func NewConversationSummaryRepo(db *gorm.DB) providers.ConversationSummaryProvider {
	return &conversationSummaryRepo{db: db}
}

// Upsert escribe (o reemplaza) el resumen de una conversación y rehace sus filas
// de join con alumnos/dispositivos, todo en una transacción. Idempotente por
// conversation_id: re-resumir una conversación pisa el resumen anterior.
func (r *conversationSummaryRepo) Upsert(ctx context.Context, s entities.ConversationSummary, studentIDs, deviceIDs []int64) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "conversation_id"}},
			DoUpdates: clause.AssignmentColumns([]string{"summary", "topic_keywords", "token_count", "updated_at"}),
		}).Create(&s).Error; err != nil {
			return fmt.Errorf("upsert summary: %w", err)
		}

		if err := tx.Exec("DELETE FROM conversation_summary_students WHERE conversation_id = ?", s.ConversationID).Error; err != nil {
			return fmt.Errorf("clear student joins: %w", err)
		}
		for _, sid := range studentIDs {
			if err := tx.Exec(
				"INSERT INTO conversation_summary_students (conversation_id, student_id) VALUES (?, ?) ON CONFLICT DO NOTHING",
				s.ConversationID, sid).Error; err != nil {
				return fmt.Errorf("insert student join: %w", err)
			}
		}

		if err := tx.Exec("DELETE FROM conversation_summary_devices WHERE conversation_id = ?", s.ConversationID).Error; err != nil {
			return fmt.Errorf("clear device joins: %w", err)
		}
		for _, did := range deviceIDs {
			if err := tx.Exec(
				"INSERT INTO conversation_summary_devices (conversation_id, device_id) VALUES (?, ?) ON CONFLICT DO NOTHING",
				s.ConversationID, did).Error; err != nil {
				return fmt.Errorf("insert device join: %w", err)
			}
		}
		return nil
	})
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
