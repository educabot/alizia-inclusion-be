package entities

import (
	"time"

	"github.com/lib/pq"
)

// ConversationSummary is the compressed summary of a conversation (1:1 with conversations).
// Written at session end and retrieved per entity on open.
type ConversationSummary struct {
	ConversationID int64          `json:"conversation_id" gorm:"primaryKey"`
	Summary        string         `json:"summary"`
	TopicKeywords  pq.StringArray `json:"topic_keywords" gorm:"type:text[]"`
	TokenCount     int            `json:"token_count"`
	UpdatedAt      time.Time      `json:"updated_at"`
}

func (ConversationSummary) TableName() string {
	return "conversation_summaries"
}
