package entrypoints

import (
	"github.com/educabot/team-ai-toolkit/web"

	"github.com/educabot/alizia-inclusion-be/src/core/usecases/management"
	"github.com/educabot/alizia-inclusion-be/src/entrypoints/middleware"
	"github.com/educabot/alizia-inclusion-be/src/entrypoints/rest"
)

func (c *ManagementContainer) HandleListTeachers(req web.Request) web.Response {
	result, err := c.ListTeachers.Execute(req.Context(), management.ListTeachersRequest{
		OrgID: middleware.OrgID(req),
	})
	if err != nil {
		return rest.HandleError(err)
	}
	return web.OK(mapUsers(result))
}
