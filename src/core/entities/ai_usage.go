package entities

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

type AIUsage struct {
	ID               int64     `json:"id" gorm:"primaryKey"`
	OrganizationID   uuid.UUID `json:"organization_id"`
	UserID           int64     `json:"user_id"`
	Mode             string    `json:"mode"`
	PromptTokens     int       `json:"prompt_tokens"`
	CompletionTokens int       `json:"completion_tokens"`
	TotalTokens      int       `json:"total_tokens"`
	// Traza por turno (HU-6, T-6.5; columnas nullable de la migración 000020).
	// ContextSnapshot guarda solo IDs (sin PII): student_id, classroom_id, device_ids…
	ConversationID  *int64         `json:"conversation_id,omitempty"`
	MessageID       *int64         `json:"message_id,omitempty"`
	Model           string         `json:"model,omitempty"`
	LatencyMs       int            `json:"latency_ms"`
	ToolCalls       int            `json:"tool_calls"`
	ContextSnapshot datatypes.JSON `json:"context_snapshot,omitempty" gorm:"type:jsonb;default:'{}'"`
	CreatedAt       time.Time      `json:"created_at"`
}

func (AIUsage) TableName() string { return "ai_usage" }
