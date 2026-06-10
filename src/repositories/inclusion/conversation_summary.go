package inclusion

import (
	"context"
	"fmt"
	"time"

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

// Upsert guarda/actualiza el resumen compactado de una conversación y revincula sus
// cross-tables (alumnos / devices). Es idempotente por conversation_id: re-cerrar la
// misma conversación reemplaza el resumen y sus vínculos en una sola transacción.
func (r *conversationSummaryRepo) Upsert(ctx context.Context, summary *entities.ConversationSummary, studentIDs, deviceIDs []int64) error {
	if summary == nil || summary.ConversationID == 0 {
		return fmt.Errorf("%w: conversation_id is required to upsert a summary", providers.ErrValidation)
	}

	summary.UpdatedAt = time.Now()

	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "conversation_id"}},
			DoUpdates: clause.AssignmentColumns([]string{"summary", "topic_keywords", "token_count", "updated_at"}),
		}).Create(summary).Error; err != nil {
			return fmt.Errorf("upsert conversation summary: %w", err)
		}

		if err := relinkSummary(tx, "conversation_summary_students", "student_id", summary.ConversationID, studentIDs); err != nil {
			return err
		}
		if err := relinkSummary(tx, "conversation_summary_devices", "device_id", summary.ConversationID, deviceIDs); err != nil {
			return err
		}
		return nil
	})
}

// relinkSummary reemplaza los vínculos de un resumen con una dimensión (alumnos o
// devices): borra los previos e inserta los actuales, dejando la tabla consistente
// con el cierre más reciente. ON CONFLICT DO NOTHING tolera ids repetidos.
func relinkSummary(tx *gorm.DB, table, fkColumn string, conversationID int64, ids []int64) error {
	if err := tx.Exec("DELETE FROM "+table+" WHERE conversation_id = ?", conversationID).Error; err != nil {
		return fmt.Errorf("clear %s: %w", table, err)
	}
	for _, id := range ids {
		stmt := "INSERT INTO " + table + " (conversation_id, " + fkColumn + ") VALUES (?, ?) ON CONFLICT DO NOTHING"
		if err := tx.Exec(stmt, conversationID, id).Error; err != nil {
			return fmt.Errorf("link %s: %w", table, err)
		}
	}
	return nil
}
