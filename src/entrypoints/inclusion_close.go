package entrypoints

import (
	"github.com/educabot/team-ai-toolkit/web"

	"github.com/educabot/alizia-inclusion-be/src/core/usecases/inclusion"
	"github.com/educabot/alizia-inclusion-be/src/entrypoints/middleware"
	"github.com/educabot/alizia-inclusion-be/src/entrypoints/rest"
)

type closeSessionBody struct {
	ConversationID int64 `json:"conversation_id"`
}

// HandleCloseSession compacts the conversation when the session closes (HU-5):
// it creates/updates the compressed summary in DB, tagged across three dimensions
// (student / topic / device kit).
func (c *InclusionContainer) HandleCloseSession(req web.Request) web.Response {
	var body closeSessionBody
	if err := req.BindJSON(&body); err != nil {
		return rest.HandleError(err)
	}

	result, err := c.CloseSession.Execute(req.Context(), inclusion.CloseSessionRequest{
		OrgID:          middleware.OrgID(req),
		UserID:         middleware.UserID(req),
		ConversationID: body.ConversationID,
	})
	if err != nil {
		return rest.HandleError(err)
	}
	return web.OK(result)
}
