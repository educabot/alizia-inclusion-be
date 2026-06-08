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

// HandleOpenSession es el Prompt 0 / router de apertura (HU-1): da la bienvenida, pregunta
// la dimensión (alumno / valija / tema), dirige la carga de contexto y recupera resúmenes
// previos de la entidad.
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
