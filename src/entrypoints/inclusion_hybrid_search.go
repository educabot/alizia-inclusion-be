package entrypoints

import (
	"github.com/educabot/team-ai-toolkit/web"

	"github.com/educabot/alizia-inclusion-be/src/core/usecases/inclusion"
	"github.com/educabot/alizia-inclusion-be/src/entrypoints/rest"
)

type hybridSearchBody struct {
	SemanticQuestion string   `json:"semantic_question"`
	Terms            []string `json:"terms"`
	ResourceID       *int64   `json:"resource_id"`
	Limit            int      `json:"limit"`
	Offset           int      `json:"offset"`
}

// HandleHybridSearch expone la búsqueda híbrida (vector + texto + conceptos) sobre el
// corpus rag_*. Sin LLM, para validar el corpus por Postman. El corpus es global
// (no se filtra por organización). Sin match → results vacío.
func (c *InclusionContainer) HandleHybridSearch(req web.Request) web.Response {
	var body hybridSearchBody
	if err := req.BindJSON(&body); err != nil {
		return rest.HandleError(err)
	}

	result, err := c.HybridSearchContent.Execute(req.Context(), inclusion.HybridSearchRequest{
		SemanticQuestion: body.SemanticQuestion,
		Terms:            body.Terms,
		ResourceID:       body.ResourceID,
		Limit:            body.Limit,
		Offset:           body.Offset,
	})
	if err != nil {
		return rest.HandleError(err)
	}
	return web.OK(result)
}
