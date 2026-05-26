package entities

import (
	"time"

	"github.com/google/uuid"
)

type AIUsage struct {
	ID               int64     `json:"id" gorm:"primaryKey"`
	OrganizationID   uuid.UUID `json:"organization_id"`
	UserID           int64     `json:"user_id"`
	Mode             string    `json:"mode"`
	PromptTokens     int       `json:"prompt_tokens"`
	CompletionTokens int       `json:"completion_tokens"`
	TotalTokens      int       `json:"total_tokens"`
	CreatedAt        time.Time `json:"created_at"`
}

func (AIUsage) TableName() string { return "ai_usage" }
