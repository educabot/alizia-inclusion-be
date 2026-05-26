package entrypoints

import (
	"strconv"

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

func (c *DashboardContainer) HandleGetAIUsage(req web.Request) web.Response {
	days := 0
	if raw := req.Query("days"); raw != "" {
		if parsed, perr := strconv.Atoi(raw); perr == nil {
			days = parsed
		}
	}

	result, err := c.GetAIUsage.Execute(req.Context(), dashboard.GetAIUsageRequest{
		OrgID: middleware.OrgID(req),
		Days:  days,
	})
	if err != nil {
		return rest.HandleError(err)
	}
	return web.OK(result)
}
