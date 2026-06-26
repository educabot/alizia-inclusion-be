package entities

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

type Conversation struct {
	ID             int64                 `json:"id" gorm:"primaryKey"`
	OrganizationID uuid.UUID             `json:"organization_id"`
	UserID         int64                 `json:"user_id"`
	StudentID      *int64                `json:"student_id,omitempty"`
	Mode           string                `json:"mode"`
	Title          string                `json:"title"`
	Messages       []ConversationMessage `json:"messages,omitempty" gorm:"foreignKey:ConversationID"`
	TimeTrackedEntity
}

type ConversationMessage struct {
	ID             int64          `json:"id" gorm:"primaryKey"`
	ConversationID int64          `json:"conversation_id"`
	Role           string         `json:"role"`
	Content        string         `json:"content"`
	Metadata       datatypes.JSON `json:"metadata" gorm:"type:jsonb;default:'{}'"`
	CreatedAt      time.Time      `json:"created_at"`
}
