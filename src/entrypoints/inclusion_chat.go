package entrypoints

import (
	"net/http"
	"strconv"
	"time"

	"github.com/educabot/team-ai-toolkit/web"

	"github.com/educabot/alizia-inclusion-be/src/core/entities"
	"github.com/educabot/alizia-inclusion-be/src/core/usecases/inclusion"
	"github.com/educabot/alizia-inclusion-be/src/entrypoints/middleware"
	"github.com/educabot/alizia-inclusion-be/src/entrypoints/rest"
)

type conversationResponse struct {
	ID        int64                     `json:"id"`
	Mode      string                    `json:"mode"`
	Messages  []conversationMsgResponse `json:"messages"`
	CreatedAt string                    `json:"created_at"`
}

type conversationMsgResponse struct {
	// ID es el id persistido del mensaje; el FE lo usa para anclar el feedback.
	ID        int64  `json:"id"`
	Role      string `json:"role"`
	Content   string `json:"content"`
	CreatedAt string `json:"created_at"`
}

func mapConversation(c entities.Conversation) conversationResponse {
	msgs := make([]conversationMsgResponse, len(c.Messages))
	for i, m := range c.Messages {
		msgs[i] = conversationMsgResponse{
			ID:        m.ID,
			Role:      m.Role,
			Content:   m.Content,
			CreatedAt: m.CreatedAt.Format(time.RFC3339),
		}
	}
	return conversationResponse{
		ID:        c.ID,
		Mode:      c.Mode,
		Messages:  msgs,
		CreatedAt: c.CreatedAt.Format(time.RFC3339),
	}
}

func mapConversations(cs []entities.Conversation) []conversationResponse {
	out := make([]conversationResponse, len(cs))
	for i := range cs {
		out[i] = mapConversation(cs[i])
	}
	return out
}

func (c *InclusionContainer) HandleGetChatHistory(req web.Request) web.Response {
	mode := req.Param("contextId")

	result, err := c.GetChatHistory.Execute(req.Context(), inclusion.GetChatHistoryRequest{
		OrgID:  middleware.OrgID(req),
		UserID: middleware.UserID(req),
		Mode:   mode,
	})
	if err != nil {
		return rest.HandleError(err)
	}
	return web.OK(mapConversations(result))
}

func (c *InclusionContainer) HandleDeleteConversation(req web.Request) web.Response {
	id, err := strconv.ParseInt(req.Param("id"), 10, 64)
	if err != nil {
		return rest.HandleError(err)
	}

	if err := c.DeleteConversation.Execute(req.Context(), inclusion.DeleteConversationRequest{
		OrgID:          middleware.OrgID(req),
		ConversationID: id,
	}); err != nil {
		return rest.HandleError(err)
	}
	return web.Response{Status: http.StatusNoContent}
}

type renameConversationBody struct {
	Title string `json:"title"`
}

func (c *InclusionContainer) HandleRenameConversation(req web.Request) web.Response {
	id, err := strconv.ParseInt(req.Param("id"), 10, 64)
	if err != nil {
		return rest.HandleError(err)
	}

	var body renameConversationBody
	if err := req.BindJSON(&body); err != nil {
		return rest.HandleError(err)
	}

	if err := c.RenameConversation.Execute(req.Context(), inclusion.RenameConversationRequest{
		OrgID:          middleware.OrgID(req),
		ConversationID: id,
		Title:          body.Title,
	}); err != nil {
		return rest.HandleError(err)
	}
	return web.Response{Status: http.StatusNoContent}
}
