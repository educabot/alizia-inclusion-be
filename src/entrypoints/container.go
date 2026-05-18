package entrypoints

import (
	"github.com/educabot/team-ai-toolkit/web"
)

type WebHandlerContainer struct {
	Auth             *AuthContainer
	Catalog          *CatalogContainer
	Inclusion        *InclusionContainer
	Management       *ManagementContainer
	Dashboard        *DashboardContainer
	AuthMiddleware   web.Interceptor
	TenantMiddleware web.Interceptor
}
