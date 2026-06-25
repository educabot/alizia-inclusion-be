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

// RecentByStudent returns the most recent summaries linked to a student, scoped to the given org.
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

// RecentByDevice returns the most recent summaries linked to a suitcase device.
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

// RecentByTopic returns the most recent summaries whose topic_keywords array contains the given keyword.
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

// Upsert saves or updates the compacted summary for a conversation and re-links its
// cross-tables (students / devices). Idempotent on conversation_id: closing the same
// conversation again replaces the summary and all its links in a single transaction.
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

// relinkSummary replaces all links for a summary along one dimension (students or devices):
// deletes existing rows then inserts the current set, keeping the table consistent with
// the latest close. ON CONFLICT DO NOTHING tolerates duplicate IDs in the input slice.
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
