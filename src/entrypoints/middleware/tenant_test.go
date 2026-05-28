package middleware_test

import (
	"testing"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

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

	require.Equal(t, 0, resp.Status)
	assert.Equal(t, orgUUID, mw.OrgID(req))
	assert.Equal(t, int64(42), mw.UserID(req))
}

func TestTenantMiddleware_RejectsMissingClaims(t *testing.T) {
	interceptor := mw.TenantMiddleware()
	req := web.NewMockRequest()

	resp := interceptor(req)

	assert.Equal(t, 401, resp.Status)
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

	assert.Equal(t, 401, resp.Status)
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

	assert.Equal(t, 401, resp.Status)
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

	assert.Equal(t, 401, resp.Status)
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

	assert.Equal(t, 401, resp.Status)
}

func TestOrgID_ReturnsNilWhenNotSet(t *testing.T) {
	req := web.NewMockRequest()
	assert.Equal(t, uuid.Nil, mw.OrgID(req))
}

func TestOrgID_ReturnsNilForWrongType(t *testing.T) {
	req := web.NewMockRequest()
	req.Values[mw.OrgIDKey] = "not-a-uuid-type"
	assert.Equal(t, uuid.Nil, mw.OrgID(req))
}

func TestUserID_ReturnsZeroWhenNotSet(t *testing.T) {
	req := web.NewMockRequest()
	assert.Equal(t, int64(0), mw.UserID(req))
}

func TestUserID_ReturnsZeroForWrongType(t *testing.T) {
	req := web.NewMockRequest()
	req.Values[mw.UserIDKey] = "not-an-int"
	assert.Equal(t, int64(0), mw.UserID(req))
}
