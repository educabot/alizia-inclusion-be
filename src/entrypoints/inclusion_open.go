package entrypoints

import (
	"github.com/educabot/team-ai-toolkit/web"

	"github.com/educabot/alizia-inclusion-be/src/core/usecases/inclusion"
	"github.com/educabot/alizia-inclusion-be/src/entrypoints/middleware"
	"github.com/educabot/alizia-inclusion-be/src/entrypoints/rest"
)

type openSessionBody struct {
	Dimension string `json:"dimension"`
	StudentID *int64 `json:"student_id"`
	DeviceID  *int64 `json:"device_id"`
	Topic     string `json:"topic"`
}

// HandleOpenSession is the Prompt 0 / opening router: greets the user, asks for the
// dimension (student / kit / topic), directs context loading, and retrieves prior summaries
// for the entity.
func (c *InclusionContainer) HandleOpenSession(req web.Request) web.Response {
	var body openSessionBody
	if err := req.BindJSON(&body); err != nil {
		return rest.HandleError(err)
	}

	result, err := c.OpenSession.Execute(req.Context(), inclusion.OpenSessionRequest{
		OrgID:     middleware.OrgID(req),
		UserID:    middleware.UserID(req),
		Dimension: body.Dimension,
		StudentID: body.StudentID,
		DeviceID:  body.DeviceID,
		Topic:     body.Topic,
	})
	if err != nil {
		return rest.HandleError(err)
	}
	return web.OK(result)
}
