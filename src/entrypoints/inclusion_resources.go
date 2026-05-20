package entrypoints

import (
	"strconv"
	"time"

	"github.com/educabot/team-ai-toolkit/web"

	"github.com/educabot/alizia-inclusion-be/src/core/entities"
	"github.com/educabot/alizia-inclusion-be/src/core/usecases/inclusion"
	"github.com/educabot/alizia-inclusion-be/src/entrypoints/rest"
)

type adaptationResourceResponse struct {
	ID           int64  `json:"id"`
	AdaptationID int64  `json:"adaptation_id"`
	Title        string `json:"title"`
	FileURL      string `json:"file_url"`
	FileType     string `json:"file_type"`
	CreatedAt    string `json:"created_at"`
}

func mapResource(r entities.AdaptationResource) adaptationResourceResponse {
	return adaptationResourceResponse{
		ID:           r.ID,
		AdaptationID: r.AdaptationID,
		Title:        r.Title,
		FileURL:      r.FileURL,
		FileType:     r.FileType,
		CreatedAt:    r.CreatedAt.Format(time.RFC3339),
	}
}

func mapResources(rs []entities.AdaptationResource) []adaptationResourceResponse {
	out := make([]adaptationResourceResponse, len(rs))
	for i, r := range rs {
		out[i] = mapResource(r)
	}
	return out
}

func (c *InclusionContainer) HandleListAdaptationResources(req web.Request) web.Response {
	id, err := strconv.ParseInt(req.Param("id"), 10, 64)
	if err != nil {
		return rest.HandleError(err)
	}

	result, err := c.ListAdaptationResources.Execute(req.Context(), inclusion.ListAdaptationResourcesRequest{
		AdaptationID: id,
	})
	if err != nil {
		return rest.HandleError(err)
	}
	return web.OK(mapResources(result))
}
