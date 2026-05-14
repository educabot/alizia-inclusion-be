package entrypoints

import (
	"strconv"

	"github.com/educabot/team-ai-toolkit/web"

	"github.com/educabot/alizia-inclusion-be/src/core/entities"
	"github.com/educabot/alizia-inclusion-be/src/core/usecases/catalog"
	"github.com/educabot/alizia-inclusion-be/src/entrypoints/middleware"
	"github.com/educabot/alizia-inclusion-be/src/entrypoints/rest"
)

type deviceResponse struct {
	ID                 int64   `json:"id"`
	RampID             int64   `json:"ramp_id"`
	Name               string  `json:"name"`
	Description        *string `json:"description,omitempty"`
	ImageURL           *string `json:"image_url,omitempty"`
	QRCode             *string `json:"qr_code,omitempty"`
	HowToUse           *string `json:"how_to_use,omitempty"`
	Recommendations    *string `json:"recommendations,omitempty"`
	Rationale          *string `json:"rationale,omitempty"`
	ClassroomBenefit   *string `json:"classroom_benefit,omitempty"`
	NeedsDescription   *string `json:"needs_description,omitempty"`
	EvaluationCriteria *string `json:"evaluation_criteria,omitempty"`
	Quantity           int     `json:"quantity"`
	SortOrder          int     `json:"sort_order"`
	RampName           string  `json:"ramp_name,omitempty"`
}

func mapDevice(d entities.Device) deviceResponse {
	resp := deviceResponse{
		ID:                 d.ID,
		RampID:             d.RampID,
		Name:               d.Name,
		Description:        d.Description,
		ImageURL:           d.ImageURL,
		QRCode:             d.QRCode,
		HowToUse:           d.HowToUse,
		Recommendations:    d.Recommendations,
		Rationale:          d.Rationale,
		ClassroomBenefit:   d.ClassroomBenefit,
		NeedsDescription:   d.NeedsDescription,
		EvaluationCriteria: d.EvaluationCriteria,
		Quantity:           d.Quantity,
		SortOrder:          d.SortOrder,
	}
	if d.Ramp != nil {
		resp.RampName = d.Ramp.Name
	}
	return resp
}

func mapDevices(ds []entities.Device) []deviceResponse {
	out := make([]deviceResponse, len(ds))
	for i, d := range ds {
		out[i] = mapDevice(d)
	}
	return out
}

func (c *CatalogContainer) HandleListDevices(req web.Request) web.Response {
	var rampID *int64
	if v := req.Query("ramp_id"); v != "" {
		id, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			return rest.HandleError(err)
		}
		rampID = &id
	}

	result, err := c.ListDevices.Execute(req.Context(), catalog.ListDevicesRequest{
		OrgID:  middleware.OrgID(req),
		RampID: rampID,
	})
	if err != nil {
		return rest.HandleError(err)
	}
	return web.OK(mapDevices(result))
}

func (c *CatalogContainer) HandleGetDevice(req web.Request) web.Response {
	id, err := strconv.ParseInt(req.Param("id"), 10, 64)
	if err != nil {
		return rest.HandleError(err)
	}

	result, err := c.GetDevice.Execute(req.Context(), catalog.GetDeviceRequest{
		OrgID:    middleware.OrgID(req),
		DeviceID: id,
	})
	if err != nil {
		return rest.HandleError(err)
	}
	return web.OK(mapDevice(*result))
}
