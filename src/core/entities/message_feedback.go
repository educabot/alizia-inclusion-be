package entities

import (
	"github.com/google/uuid"
)

// MessageFeedback es el feedback (manito arriba/abajo + comentario opcional) que
// un usuario deja sobre un mensaje del asistente. Uso interno: sirve para revisar
// los errores de Alizia. Guardamos ConversationID además del mensaje para poder
// reconstruir el hilo completo y entender el contexto del error.
type MessageFeedback struct {
	ID                    int64     `json:"id" gorm:"primaryKey"`
	ConversationMessageID int64     `json:"conversation_message_id"`
	ConversationID        int64     `json:"conversation_id"`
	OrganizationID        uuid.UUID `json:"organization_id"`
	UserID                int64     `json:"user_id"`
	Rating                string    `json:"rating"` // like | dislike
	Comment               string    `json:"comment"`
	TimeTrackedEntity
}

func (MessageFeedback) TableName() string {
	return "message_feedback"
}
