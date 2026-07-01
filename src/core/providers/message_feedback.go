package providers

import (
	"context"

	"github.com/google/uuid"

	"github.com/educabot/alizia-inclusion-be/src/core/entities"
)

// MessageFeedbackReview es una fila de la vista de revisión interna: el feedback
// más el contenido del mensaje comentado y la pregunta del usuario que lo generó,
// para entender el contexto del error sin abrir la conversación entera.
type MessageFeedbackReview struct {
	entities.MessageFeedback
	MessageContent      string `json:"message_content" gorm:"column:message_content"`
	PreviousUserMessage string `json:"previous_user_message" gorm:"column:previous_user_message"`
}

type MessageFeedbackProvider interface {
	// MessageContext devuelve el conversation_id y la org dueña del mensaje del
	// asistente, para validar pertenencia antes de guardar. ErrNotFound si no existe.
	MessageContext(ctx context.Context, messageID int64) (conversationID int64, orgID uuid.UUID, err error)
	// Upsert crea o actualiza el feedback del usuario para el mensaje (unique por
	// conversation_message_id + user_id).
	Upsert(ctx context.Context, fb *entities.MessageFeedback) error
	// Delete borra el feedback del usuario para el mensaje (toggle-off).
	Delete(ctx context.Context, orgID uuid.UUID, messageID, userID int64) error
	// List devuelve los feedbacks de la organización para revisión interna,
	// opcionalmente filtrados por rating ("" = todos).
	List(ctx context.Context, orgID uuid.UUID, rating string) ([]MessageFeedbackReview, error)
}
