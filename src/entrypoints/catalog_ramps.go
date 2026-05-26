package entrypoints

import (
	"strconv"

	"github.com/educabot/team-ai-toolkit/web"

	"github.com/educabot/alizia-inclusion-be/src/core/entities"
	"github.com/educabot/alizia-inclusion-be/src/core/usecases/catalog"
	"github.com/educabot/alizia-inclusion-be/src/entrypoints/middleware"
	"github.com/educabot/alizia-inclusion-be/src/entrypoints/rest"
)

type rampResponse struct {
	ID               int64            `json:"id"`
	Name             string           `json:"name"`
	Description      *string          `json:"description,omitempty"`
	ShortDescription *string          `json:"short_description,omitempty"`
	VideoURL         *string          `json:"video_url,omitempty"`
	SortOrder        int              `json:"sort_order"`
	Devices          []deviceResponse `json:"devices"`
}

func mapRamp(r entities.Ramp) rampResponse {
	return rampResponse{
		ID:               r.ID,
		Name:             r.Name,
		Description:      r.Description,
		ShortDescription: r.ShortDescription,
		VideoURL:         r.VideoURL,
		SortOrder:        r.SortOrder,
		Devices:          mapDevices(r.Devices),
	}
}

func mapRamps(rs []entities.Ramp) []rampResponse {
	out := make([]rampResponse, len(rs))
	for i := range rs {
		out[i] = mapRamp(rs[i])
	}
	return out
}

func (c *CatalogContainer) HandleListRamps(req web.Request) web.Response {
	result, err := c.ListRamps.Execute(req.Context(), catalog.ListRampsRequest{
		OrgID: middleware.OrgID(req),
	})
	if err != nil {
		return rest.HandleError(err)
	}
	return web.OK(mapRamps(result))
}

func (c *CatalogContainer) HandleGetRamp(req web.Request) web.Response {
	id, err := strconv.ParseInt(req.Param("id"), 10, 64)
	if err != nil {
		return rest.HandleError(err)
	}

	result, err := c.GetRamp.Execute(req.Context(), catalog.GetRampRequest{
		OrgID:  middleware.OrgID(req),
		RampID: id,
	})
	if err != nil {
		return rest.HandleError(err)
	}
	return web.OK(mapRamp(*result))
}
