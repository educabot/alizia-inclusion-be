package middleware

import (
	"github.com/educabot/team-ai-toolkit/config"
	"github.com/educabot/team-ai-toolkit/tokens"
	"github.com/educabot/team-ai-toolkit/web"
)

// RS256AuthMiddleware validates Bearer tokens issued by the external auth-service (RS256).
//
// Validation is delegated to the shared tokens.AuthServiceMiddleware (toolkit v1.9.0+) so
// every consumer of auth-service uses the SAME implementation: it checks RS256 +
// iss=auth-service + aud=educabot-api + exp, parses typed claims, and projects the org_uuid
// into tokens.Claims (Audience) so the existing TenantMiddleware keeps reading it unchanged.
//
// This service is multi-tenant: ExpectedOrgUUID is empty, so no organization is rejected —
// each request is scoped downstream by the org_uuid carried in its own token. In config.Test
// the toolkit injects a fixed mock identity (org 00000000-0000-0000-0000-000000000001).
func RS256AuthMiddleware(publicKeyPEM string, env config.Environment) web.Interceptor {
	return tokens.AuthServiceMiddleware(tokens.AuthServiceConfig{
		PublicKeyPEM:    publicKeyPEM,
		ExpectedOrgUUID: "",
		Env:             env,
	})
}
