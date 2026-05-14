package entrypoints

import (
	"github.com/educabot/team-ai-toolkit/web"
)

type WebHandlerContainer struct {
	Catalog          *CatalogContainer
	Inclusion        *InclusionContainer
	AuthMiddleware   web.Interceptor
	TenantMiddleware web.Interceptor
}
