package middleware_test

import (
	"testing"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"

	"github.com/educabot/team-ai-toolkit/tokens"
	"github.com/educabot/team-ai-toolkit/web"

	mw "github.com/educabot/alizia-inclusion-be/src/entrypoints/middleware"
)

func TestTenantMiddleware_SetsOrgIDAndUserIDFromValidClaims(t *testing.T) {
	interceptor := mw.TenantMiddleware()
	orgUUID := uuid.New()
	req := web.NewMockRequest()
	req.Values[tokens.ClaimsKey] = &tokens.Claims{
		ID: "42",
		RegisteredClaims: jwt.RegisteredClaims{
			Audience: jwt.ClaimStrings{orgUUID.String()},
		},
	}

	resp := interceptor(req)

	if resp.Status != 0 {
		t.Fatalf("expected pass-through (status 0), got %d", resp.Status)
	}

	gotOrg := mw.OrgID(req)
	if gotOrg != orgUUID {
		t.Errorf("expected org_id %s, got %s", orgUUID, gotOrg)
	}

	gotUser := mw.UserID(req)
	if gotUser != 42 {
		t.Errorf("expected user_id 42, got %d", gotUser)
	}
}

func TestTenantMiddleware_RejectsMissingClaims(t *testing.T) {
	interceptor := mw.TenantMiddleware()
	req := web.NewMockRequest()

	resp := interceptor(req)

	if resp.Status != 401 {
		t.Errorf("expected 401, got %d", resp.Status)
	}
}

func TestTenantMiddleware_RejectsEmptyAudience(t *testing.T) {
	interceptor := mw.TenantMiddleware()
	req := web.NewMockRequest()
	req.Values[tokens.ClaimsKey] = &tokens.Claims{
		ID: "1",
		RegisteredClaims: jwt.RegisteredClaims{
			Audience: jwt.ClaimStrings{},
		},
	}

	resp := interceptor(req)

	if resp.Status != 401 {
		t.Errorf("expected 401 for empty audience, got %d", resp.Status)
	}
}

func TestTenantMiddleware_RejectsInvalidUUIDInAudience(t *testing.T) {
	interceptor := mw.TenantMiddleware()
	req := web.NewMockRequest()
	req.Values[tokens.ClaimsKey] = &tokens.Claims{
		ID: "1",
		RegisteredClaims: jwt.RegisteredClaims{
			Audience: jwt.ClaimStrings{"not-a-uuid"},
		},
	}

	resp := interceptor(req)

	if resp.Status != 401 {
		t.Errorf("expected 401 for invalid UUID, got %d", resp.Status)
	}
}

func TestTenantMiddleware_RejectsNonNumericUserID(t *testing.T) {
	interceptor := mw.TenantMiddleware()
	orgUUID := uuid.New()
	req := web.NewMockRequest()
	req.Values[tokens.ClaimsKey] = &tokens.Claims{
		ID: "abc",
		RegisteredClaims: jwt.RegisteredClaims{
			Audience: jwt.ClaimStrings{orgUUID.String()},
		},
	}

	resp := interceptor(req)

	if resp.Status != 401 {
		t.Errorf("expected 401 for non-numeric user_id, got %d", resp.Status)
	}
}

func TestTenantMiddleware_RejectsZeroUserID(t *testing.T) {
	interceptor := mw.TenantMiddleware()
	orgUUID := uuid.New()
	req := web.NewMockRequest()
	req.Values[tokens.ClaimsKey] = &tokens.Claims{
		ID: "0",
		RegisteredClaims: jwt.RegisteredClaims{
			Audience: jwt.ClaimStrings{orgUUID.String()},
		},
	}

	resp := interceptor(req)

	if resp.Status != 401 {
		t.Errorf("expected 401 for zero user_id, got %d", resp.Status)
	}
}

func TestOrgID_ReturnsNilWhenNotSet(t *testing.T) {
	req := web.NewMockRequest()
	if got := mw.OrgID(req); got != uuid.Nil {
		t.Errorf("expected uuid.Nil, got %s", got)
	}
}

func TestOrgID_ReturnsNilForWrongType(t *testing.T) {
	req := web.NewMockRequest()
	req.Values[mw.OrgIDKey] = "not-a-uuid-type"
	if got := mw.OrgID(req); got != uuid.Nil {
		t.Errorf("expected uuid.Nil for wrong type, got %s", got)
	}
}

func TestUserID_ReturnsZeroWhenNotSet(t *testing.T) {
	req := web.NewMockRequest()
	if got := mw.UserID(req); got != 0 {
		t.Errorf("expected 0, got %d", got)
	}
}

func TestUserID_ReturnsZeroForWrongType(t *testing.T) {
	req := web.NewMockRequest()
	req.Values[mw.UserIDKey] = "not-an-int"
	if got := mw.UserID(req); got != 0 {
		t.Errorf("expected 0 for wrong type, got %d", got)
	}
}
