package middleware

import (
	"crypto/rsa"
	"net/http"
	"strconv"
	"strings"

	"github.com/golang-jwt/jwt/v5"

	bcfg "github.com/educabot/team-ai-toolkit/config"
	"github.com/educabot/team-ai-toolkit/tokens"
	"github.com/educabot/team-ai-toolkit/web"
)

// authServiceClaims mirrors the AccessTokenClaims struct issued by the auth-service.
type authServiceClaims struct {
	UserID  int64    `json:"sub"`
	OrgID   int64    `json:"org_id"`
	OrgUUID string   `json:"org_uuid"`
	Roles   []string `json:"roles"`
	Email   string   `json:"email"`
	Name    string   `json:"name"`
	jwt.RegisteredClaims
}

// RS256AuthMiddleware validates Bearer tokens signed with RS256 by the auth-service.
// The public key is parsed once at construction time.
// In test environments it falls back to a mock that injects a fixed test user.
func RS256AuthMiddleware(publicKeyPEM string, env bcfg.Environment) web.Interceptor {
	if env == bcfg.Test {
		return mockAuthMiddleware()
	}

	pubKey, err := jwt.ParseRSAPublicKeyFromPEM([]byte(publicKeyPEM))
	if err != nil {
		// Fail fast: a misconfigured key should prevent the server from accepting requests.
		panic("auth middleware: invalid RSA public key: " + err.Error())
	}

	return rs256Middleware(pubKey)
}

func rs256Middleware(pubKey *rsa.PublicKey) web.Interceptor {
	keyFunc := func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return pubKey, nil
	}

	return func(req web.Request) web.Response {
		authHeader := req.Header("Authorization")
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "bearer") {
			return web.Err(http.StatusUnauthorized, "unauthorized", "missing bearer token")
		}
		rawToken := parts[1]

		var svcClaims authServiceClaims
		_, err := jwt.ParseWithClaims(rawToken, &svcClaims, keyFunc,
			jwt.WithAudience("educabot-api"),
			jwt.WithIssuer("auth-service"),
		)
		if err != nil {
			return web.Err(http.StatusUnauthorized, "unauthorized", "invalid token")
		}

		claims := &tokens.Claims{
			ID:    strconv.FormatInt(svcClaims.UserID, 10),
			Name:  svcClaims.Name,
			Email: svcClaims.Email,
			Roles: svcClaims.Roles,
			RegisteredClaims: jwt.RegisteredClaims{
				// Map OrgUUID into Audience so TenantMiddleware can read it unchanged.
				Audience: jwt.ClaimStrings{svcClaims.OrgUUID},
			},
		}

		req.Set(tokens.ClaimsKey, claims)
		req.Set(tokens.IDKey, claims.ID)
		req.Set(tokens.EmailKey, claims.Email)
		req.Set(tokens.TokenKey, rawToken)

		return web.Response{}
	}
}

// mockAuthMiddleware injects a fixed test identity, matching the toolkit's test behavior.
func mockAuthMiddleware() web.Interceptor {
	const testUserID = "1"
	const testOrgUUID = "00000000-0000-0000-0000-000000000001"

	claims := &tokens.Claims{
		ID:    testUserID,
		Name:  "Test User",
		Email: "test@educabot.com",
		Roles: []string{"teacher"},
		RegisteredClaims: jwt.RegisteredClaims{
			Audience: jwt.ClaimStrings{testOrgUUID},
		},
	}

	return func(req web.Request) web.Response {
		req.Set(tokens.ClaimsKey, claims)
		req.Set(tokens.IDKey, claims.ID)
		req.Set(tokens.EmailKey, claims.Email)
		req.Set(tokens.TokenKey, "mock-token")
		return web.Response{}
	}
}
