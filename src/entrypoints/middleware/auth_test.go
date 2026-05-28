package middleware_test

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"strconv"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	bcfg "github.com/educabot/team-ai-toolkit/config"
	"github.com/educabot/team-ai-toolkit/tokens"
	"github.com/educabot/team-ai-toolkit/web"

	mw "github.com/educabot/alizia-inclusion-be/src/entrypoints/middleware"
)

func generateRSAKeyPair(t *testing.T) (privKey *rsa.PrivateKey, publicPEM string) {
	t.Helper()
	privKey, err := rsa.GenerateKey(rand.Reader, 2048)
	require.NoError(t, err)
	pubBytes, err := x509.MarshalPKIXPublicKey(&privKey.PublicKey)
	require.NoError(t, err)
	pubPEM := pem.EncodeToMemory(&pem.Block{Type: "PUBLIC KEY", Bytes: pubBytes})
	return privKey, string(pubPEM)
}

func signAuthServiceToken(t *testing.T, privKey *rsa.PrivateKey, userID int64, orgUUID string, expired bool) string {
	t.Helper()
	exp := time.Now().Add(time.Hour)
	if expired {
		exp = time.Now().Add(-time.Hour)
	}

	claims := jwt.MapClaims{
		"sub":      userID,
		"org_uuid": orgUUID,
		"roles":    []string{"teacher"},
		"email":    "test@educabot.com",
		"name":     "Test Teacher",
		"aud":      "educabot-api",
		"iss":      "auth-service",
		"exp":      exp.Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	signed, err := token.SignedString(privKey)
	require.NoError(t, err)
	return signed
}

func TestRS256AuthMiddleware_TestEnv_InjectsMockClaims(t *testing.T) {
	interceptor := mw.RS256AuthMiddleware("", bcfg.Test)
	req := web.NewMockRequest()

	resp := interceptor(req)

	assert.Equal(t, 0, resp.Status)

	claims := tokens.GetClaims(req)
	require.NotNil(t, claims)
	assert.Equal(t, "1", claims.ID)
	assert.Equal(t, "test@educabot.com", claims.Email)
}

func TestRS256AuthMiddleware_Prod_AcceptsValidRS256Token(t *testing.T) {
	privKey, pubPEM := generateRSAKeyPair(t)
	interceptor := mw.RS256AuthMiddleware(pubPEM, bcfg.Production)
	orgUUID := "00000000-0000-0000-0000-000000000001"

	tokenStr := signAuthServiceToken(t, privKey, 42, orgUUID, false)
	req := web.NewMockRequest()
	req.Headers["Authorization"] = "Bearer " + tokenStr

	resp := interceptor(req)

	assert.Equal(t, 0, resp.Status)

	claims := tokens.GetClaims(req)
	require.NotNil(t, claims)
	userID, err := strconv.ParseInt(claims.ID, 10, 64)
	require.NoError(t, err)
	assert.Equal(t, int64(42), userID)
	assert.NotEmpty(t, claims.Audience)
	assert.Equal(t, orgUUID, claims.Audience[0])
}

func TestRS256AuthMiddleware_Prod_RejectsMissingBearerToken(t *testing.T) {
	_, pubPEM := generateRSAKeyPair(t)
	interceptor := mw.RS256AuthMiddleware(pubPEM, bcfg.Production)
	req := web.NewMockRequest()

	resp := interceptor(req)

	assert.Equal(t, 401, resp.Status)
}

func TestRS256AuthMiddleware_Prod_RejectsMalformedAuthorizationHeader(t *testing.T) {
	_, pubPEM := generateRSAKeyPair(t)
	interceptor := mw.RS256AuthMiddleware(pubPEM, bcfg.Production)
	req := web.NewMockRequest()
	req.Headers["Authorization"] = "Basic dGVzdDp0ZXN0"

	resp := interceptor(req)

	assert.Equal(t, 401, resp.Status)
}

func TestRS256AuthMiddleware_Prod_RejectsExpiredToken(t *testing.T) {
	privKey, pubPEM := generateRSAKeyPair(t)
	interceptor := mw.RS256AuthMiddleware(pubPEM, bcfg.Production)
	orgUUID := "00000000-0000-0000-0000-000000000001"

	tokenStr := signAuthServiceToken(t, privKey, 42, orgUUID, true)
	req := web.NewMockRequest()
	req.Headers["Authorization"] = "Bearer " + tokenStr

	resp := interceptor(req)

	assert.Equal(t, 401, resp.Status)
}

func TestRS256AuthMiddleware_Prod_RejectsTokenSignedWithWrongKey(t *testing.T) {
	_, pubPEM := generateRSAKeyPair(t)
	interceptor := mw.RS256AuthMiddleware(pubPEM, bcfg.Production)
	orgUUID := "00000000-0000-0000-0000-000000000001"

	wrongKey, err := rsa.GenerateKey(rand.Reader, 2048)
	require.NoError(t, err)
	tokenStr := signAuthServiceToken(t, wrongKey, 42, orgUUID, false)
	req := web.NewMockRequest()
	req.Headers["Authorization"] = "Bearer " + tokenStr

	resp := interceptor(req)

	assert.Equal(t, 401, resp.Status)
}

func TestRS256AuthMiddleware_PanicsOnInvalidKey(t *testing.T) {
	assert.Panics(t, func() {
		mw.RS256AuthMiddleware("not-a-valid-pem", bcfg.Production)
	})
}
