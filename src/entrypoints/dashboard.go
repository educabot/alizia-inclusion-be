package entrypoints

import (
	"github.com/educabot/team-ai-toolkit/web"

	"github.com/educabot/alizia-inclusion-be/src/core/usecases/dashboard"
	"github.com/educabot/alizia-inclusion-be/src/entrypoints/middleware"
	"github.com/educabot/alizia-inclusion-be/src/entrypoints/rest"
)

func (c *DashboardContainer) HandleGetMetrics(req web.Request) web.Response {
	result, err := c.GetMetrics.Execute(req.Context(), dashboard.GetMetricsRequest{
		OrgID: middleware.OrgID(req),
	})
	if err != nil {
		return rest.HandleError(err)
	}
	return web.OK(result)
}
