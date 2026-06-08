package entities

import (
	"time"

	"github.com/lib/pq"
)

// ConversationSummary es el resumen comprimido de una conversación (1:1 con conversations).
// Se escribe al fin de sesión (HU-5) y se recupera por entidad al abrir (HU-1).
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
