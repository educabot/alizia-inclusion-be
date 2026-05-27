package entrypoints

import (
	"strconv"
	"time"

	"github.com/educabot/team-ai-toolkit/web"

	"github.com/educabot/alizia-inclusion-be/src/core/entities"
	"github.com/educabot/alizia-inclusion-be/src/core/usecases/catalog"
	"github.com/educabot/alizia-inclusion-be/src/entrypoints/middleware"
	"github.com/educabot/alizia-inclusion-be/src/entrypoints/rest"
)

type deviceDownloadResponse struct {
	ID        int64  `json:"id"`
	Title     string `json:"title"`
	FileURL   string `json:"file_url"`
	FileType  string `json:"file_type"`
	CreatedAt string `json:"created_at"`
}

type deviceResponse struct {
	ID                 int64                    `json:"id"`
	RampID             int64                    `json:"ramp_id"`
	Name               string                   `json:"name"`
	Description        *string                  `json:"description,omitempty"`
	ImageURL           *string                  `json:"image_url,omitempty"`
	QRCode             *string                  `json:"qr_code,omitempty"`
	HowToUse           *string                  `json:"how_to_use,omitempty"`
	Recommendations    *string                  `json:"recommendations,omitempty"`
	Rationale          *string                  `json:"rationale,omitempty"`
	ClassroomBenefit   *string                  `json:"classroom_benefit,omitempty"`
	NeedsDescription   *string                  `json:"needs_description,omitempty"`
	UsefulWhen         *string                  `json:"useful_when,omitempty"`
	EvaluationCriteria *string                  `json:"evaluation_criteria,omitempty"`
	Quantity           int                      `json:"quantity"`
	SortOrder          int                      `json:"sort_order"`
	RampName           string                   `json:"ramp_name,omitempty"`
	Downloads          []deviceDownloadResponse `json:"downloads,omitempty"`
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
		UsefulWhen:         d.UsefulWhen,
		EvaluationCriteria: d.EvaluationCriteria,
		Quantity:           d.Quantity,
		SortOrder:          d.SortOrder,
	}
	if d.Ramp != nil {
		resp.RampName = d.Ramp.Name
	}
	for _, r := range d.Resources {
		resp.Downloads = append(resp.Downloads, deviceDownloadResponse{
			ID:        r.ID,
			Title:     r.Title,
			FileURL:   r.FileURL,
			FileType:  r.FileType,
			CreatedAt: r.CreatedAt.Format(time.RFC3339),
		})
	}
	return resp
}

func mapDevices(ds []entities.Device) []deviceResponse {
	out := make([]deviceResponse, len(ds))
	for i := range ds {
		out[i] = mapDevice(ds[i])
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
