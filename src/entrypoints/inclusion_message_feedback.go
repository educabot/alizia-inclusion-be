package entrypoints

import (
	"net/http"
	"strconv"
	"time"

	"github.com/educabot/team-ai-toolkit/web"

	"github.com/educabot/alizia-inclusion-be/src/core/entities"
	"github.com/educabot/alizia-inclusion-be/src/core/providers"
	"github.com/educabot/alizia-inclusion-be/src/core/usecases/inclusion"
	"github.com/educabot/alizia-inclusion-be/src/entrypoints/middleware"
	"github.com/educabot/alizia-inclusion-be/src/entrypoints/rest"
)

type messageFeedbackResponse struct {
	ID                    int64  `json:"id"`
	ConversationMessageID int64  `json:"conversation_message_id"`
	ConversationID        int64  `json:"conversation_id"`
	Rating                string `json:"rating"`
	Comment               string `json:"comment"`
	CreatedAt             string `json:"created_at"`
	UpdatedAt             string `json:"updated_at"`
}

func mapMessageFeedback(f entities.MessageFeedback) messageFeedbackResponse {
	return messageFeedbackResponse{
		ID:                    f.ID,
		ConversationMessageID: f.ConversationMessageID,
		ConversationID:        f.ConversationID,
		Rating:                f.Rating,
		Comment:               f.Comment,
		CreatedAt:             f.CreatedAt.Format(time.RFC3339),
		UpdatedAt:             f.UpdatedAt.Format(time.RFC3339),
	}
}

// messageFeedbackReviewResponse es la fila de la vista de revisión interna: el
// feedback + el mensaje comentado + la pregunta previa del usuario, para entender
// el contexto del error.
type messageFeedbackReviewResponse struct {
	messageFeedbackResponse
	UserID              int64  `json:"user_id"`
	MessageContent      string `json:"message_content"`
	PreviousUserMessage string `json:"previous_user_message"`
}

func mapMessageFeedbackReview(r providers.MessageFeedbackReview) messageFeedbackReviewResponse {
	return messageFeedbackReviewResponse{
		messageFeedbackResponse: mapMessageFeedback(r.MessageFeedback),
		UserID:                  r.UserID,
		MessageContent:          r.MessageContent,
		PreviousUserMessage:     r.PreviousUserMessage,
	}
}

type submitMessageFeedbackBody struct {
	Rating  string `json:"rating"`
	Comment string `json:"comment"`
}

func (c *InclusionContainer) HandleSubmitMessageFeedback(req web.Request) web.Response {
	messageID, err := strconv.ParseInt(req.Param("messageId"), 10, 64)
	if err != nil {
		return rest.HandleError(err)
	}

	var body submitMessageFeedbackBody
	if err := req.BindJSON(&body); err != nil {
		return rest.HandleError(err)
	}

	result, err := c.SubmitMessageFeedback.Execute(req.Context(), inclusion.SubmitMessageFeedbackRequest{
		OrgID:     middleware.OrgID(req),
		UserID:    middleware.UserID(req),
		MessageID: messageID,
		Rating:    body.Rating,
		Comment:   body.Comment,
	})
	if err != nil {
		return rest.HandleError(err)
	}
	return web.OK(mapMessageFeedback(*result))
}

func (c *InclusionContainer) HandleDeleteMessageFeedback(req web.Request) web.Response {
	messageID, err := strconv.ParseInt(req.Param("messageId"), 10, 64)
	if err != nil {
		return rest.HandleError(err)
	}

	if err := c.DeleteMessageFeedback.Execute(req.Context(), inclusion.DeleteMessageFeedbackRequest{
		OrgID:     middleware.OrgID(req),
		UserID:    middleware.UserID(req),
		MessageID: messageID,
	}); err != nil {
		return rest.HandleError(err)
	}
	return web.Response{Status: http.StatusNoContent}
}

// HandleListMessageFeedback lista los feedbacks de la organización para revisión
// interna. TODO: gatear por rol (director/ministerio/integradora) una vez que el
// claim de rol viaje confiable en el JWT — hoy va autenticado + scope por org.
func (c *InclusionContainer) HandleListMessageFeedback(req web.Request) web.Response {
	result, err := c.ListMessageFeedback.Execute(req.Context(), inclusion.ListMessageFeedbackRequest{
		OrgID:  middleware.OrgID(req),
		Rating: req.Query("rating"),
	})
	if err != nil {
		return rest.HandleError(err)
	}

	out := make([]messageFeedbackReviewResponse, len(result))
	for i := range result {
		out[i] = mapMessageFeedbackReview(result[i])
	}
	return web.OK(out)
}
