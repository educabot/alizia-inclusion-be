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

	bcfg "github.com/educabot/team-ai-toolkit/config"
	"github.com/educabot/team-ai-toolkit/tokens"
	"github.com/educabot/team-ai-toolkit/web"

	mw "github.com/educabot/alizia-inclusion-be/src/entrypoints/middleware"
)

func generateRSAKeyPair(t *testing.T) (privKey *rsa.PrivateKey, publicPEM string) {
	t.Helper()
	privKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatalf("generate RSA key: %v", err)
	}
	pubBytes, err := x509.MarshalPKIXPublicKey(&privKey.PublicKey)
	if err != nil {
		t.Fatalf("marshal public key: %v", err)
	}
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
	if err != nil {
		t.Fatalf("sign token: %v", err)
	}
	return signed
}

func TestRS256AuthMiddleware_TestEnv_InjectsMockClaims(t *testing.T) {
	interceptor := mw.RS256AuthMiddleware("", bcfg.Test)
	req := web.NewMockRequest()

	resp := interceptor(req)

	if resp.Status != 0 {
		t.Fatalf("expected pass-through, got status %d", resp.Status)
	}

	claims := tokens.GetClaims(req)
	if claims == nil {
		t.Fatal("expected claims to be set")
	}
	if claims.ID != "1" {
		t.Errorf("expected mock user_id '1', got %q", claims.ID)
	}
	if claims.Email != "test@educabot.com" {
		t.Errorf("expected mock email 'test@educabot.com', got %q", claims.Email)
	}
}

func TestRS256AuthMiddleware_Prod_AcceptsValidRS256Token(t *testing.T) {
	privKey, pubPEM := generateRSAKeyPair(t)
	interceptor := mw.RS256AuthMiddleware(pubPEM, bcfg.Production)
	orgUUID := "00000000-0000-0000-0000-000000000001"

	tokenStr := signAuthServiceToken(t, privKey, 42, orgUUID, false)
	req := web.NewMockRequest()
	req.Headers["Authorization"] = "Bearer " + tokenStr

	resp := interceptor(req)

	if resp.Status != 0 {
		t.Fatalf("expected pass-through, got status %d", resp.Status)
	}

	claims := tokens.GetClaims(req)
	if claims == nil {
		t.Fatal("expected claims to be set")
	}
	userID, err := strconv.ParseInt(claims.ID, 10, 64)
	if err != nil {
		t.Fatalf("parse user_id: %v", err)
	}
	if userID != 42 {
		t.Errorf("expected user_id 42, got %d", userID)
	}
	if len(claims.Audience) == 0 || claims.Audience[0] != orgUUID {
		t.Errorf("expected audience %q, got %v", orgUUID, claims.Audience)
	}
}

func TestRS256AuthMiddleware_Prod_RejectsMissingBearerToken(t *testing.T) {
	_, pubPEM := generateRSAKeyPair(t)
	interceptor := mw.RS256AuthMiddleware(pubPEM, bcfg.Production)
	req := web.NewMockRequest()

	resp := interceptor(req)

	if resp.Status != 401 {
		t.Errorf("expected 401, got %d", resp.Status)
	}
}

func TestRS256AuthMiddleware_Prod_RejectsMalformedAuthorizationHeader(t *testing.T) {
	_, pubPEM := generateRSAKeyPair(t)
	interceptor := mw.RS256AuthMiddleware(pubPEM, bcfg.Production)
	req := web.NewMockRequest()
	req.Headers["Authorization"] = "Basic dGVzdDp0ZXN0"

	resp := interceptor(req)

	if resp.Status != 401 {
		t.Errorf("expected 401 for non-bearer, got %d", resp.Status)
	}
}

func TestRS256AuthMiddleware_Prod_RejectsExpiredToken(t *testing.T) {
	privKey, pubPEM := generateRSAKeyPair(t)
	interceptor := mw.RS256AuthMiddleware(pubPEM, bcfg.Production)
	orgUUID := "00000000-0000-0000-0000-000000000001"

	tokenStr := signAuthServiceToken(t, privKey, 42, orgUUID, true)
	req := web.NewMockRequest()
	req.Headers["Authorization"] = "Bearer " + tokenStr

	resp := interceptor(req)

	if resp.Status != 401 {
		t.Errorf("expected 401 for expired token, got %d", resp.Status)
	}
}

func TestRS256AuthMiddleware_Prod_RejectsTokenSignedWithWrongKey(t *testing.T) {
	_, pubPEM := generateRSAKeyPair(t)
	interceptor := mw.RS256AuthMiddleware(pubPEM, bcfg.Production)
	orgUUID := "00000000-0000-0000-0000-000000000001"

	wrongKey, _ := rsa.GenerateKey(rand.Reader, 2048)
	tokenStr := signAuthServiceToken(t, wrongKey, 42, orgUUID, false)
	req := web.NewMockRequest()
	req.Headers["Authorization"] = "Bearer " + tokenStr

	resp := interceptor(req)

	if resp.Status != 401 {
		t.Errorf("expected 401 for wrong signing key, got %d", resp.Status)
	}
}

func TestRS256AuthMiddleware_PanicsOnInvalidKey(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("expected panic for invalid PEM key")
		}
	}()
	mw.RS256AuthMiddleware("not-a-valid-pem", bcfg.Production)
}
