package entrypoints

import (
	"github.com/educabot/team-ai-toolkit/web"

	"github.com/educabot/alizia-inclusion-be/src/core/usecases/inclusion"
	"github.com/educabot/alizia-inclusion-be/src/entrypoints/middleware"
	"github.com/educabot/alizia-inclusion-be/src/entrypoints/rest"
)

type searchContentBody struct {
	Query string `json:"query"`
	Limit int    `json:"limit"`
}

// HandleSearchContent expone el RAG de contenido pedagógico (HU-3): busca por
// keyword/full-text en el corpus y devuelve los fragmentos más relevantes. Sin
// LLM, para validar el corpus directamente. Sin match → results vacío.
func (c *InclusionContainer) HandleSearchContent(req web.Request) web.Response {
	var body searchContentBody
	if err := req.BindJSON(&body); err != nil {
		return rest.HandleError(err)
	}

	result, err := c.SearchPedagogicalContent.Execute(req.Context(), inclusion.SearchContentRequest{
		OrgID: middleware.OrgID(req),
		Query: body.Query,
		Limit: body.Limit,
	})
	if err != nil {
		return rest.HandleError(err)
	}
	return web.OK(result)
}
