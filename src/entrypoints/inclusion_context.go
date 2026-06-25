package entrypoints

import (
	"github.com/educabot/team-ai-toolkit/web"

	"github.com/educabot/alizia-inclusion-be/src/core/usecases/inclusion"
	"github.com/educabot/alizia-inclusion-be/src/entrypoints/middleware"
	"github.com/educabot/alizia-inclusion-be/src/entrypoints/rest"
)

type buildContextBody struct {
	Dimension string `json:"dimension"`
	StudentID *int64 `json:"student_id"`
	Topic     string `json:"topic"`
}

// HandleBuildContext expone el Context Assembler (HU-2): arma el contexto tipado
// para una dimensión (alumno / valija / tema) y lo devuelve. Sirve para validar
// qué contexto llega al prompt (perfil + situaciones + PPI + diagnósticos +
// adaptaciones previas), degradando con elegancia cuando faltan datos.
func (c *InclusionContainer) HandleBuildContext(req web.Request) web.Response {
	var body buildContextBody
	if err := req.BindJSON(&body); err != nil {
		return rest.HandleError(err)
	}

	result, err := c.BuildPromptContext.Execute(req.Context(), inclusion.BuildContextRequest{
		OrgID:     middleware.OrgID(req),
		UserID:    middleware.UserID(req),
		Dimension: body.Dimension,
		StudentID: body.StudentID,
		Topic:     body.Topic,
	})
	if err != nil {
		return rest.HandleError(err)
	}
	return web.OK(result)
}
