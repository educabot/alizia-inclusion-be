package entrypoints

import (
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
	Role      string `json:"role"`
	Content   string `json:"content"`
	CreatedAt string `json:"created_at"`
}

func mapConversation(c entities.Conversation) conversationResponse {
	msgs := make([]conversationMsgResponse, len(c.Messages))
	for i, m := range c.Messages {
		msgs[i] = conversationMsgResponse{
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

// HandleGetConversation returns a single conversation with its messages, scoped to
// the org. Used to resume the conversation that originated a saved resource
// (adaptation.source_conversation_id).
func (c *InclusionContainer) HandleGetConversation(req web.Request) web.Response {
	id, err := strconv.ParseInt(req.Param("id"), 10, 64)
	if err != nil {
		return rest.HandleError(err)
	}

	result, err := c.GetConversation.Execute(req.Context(), inclusion.GetConversationRequest{
		OrgID:          middleware.OrgID(req),
		ConversationID: id,
	})
	if err != nil {
		return rest.HandleError(err)
	}
	return web.OK(mapConversation(*result))
}
