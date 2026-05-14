package middleware

import (
	"net/http"
	"strconv"

	"github.com/google/uuid"
	"github.com/educabot/team-ai-toolkit/tokens"
	"github.com/educabot/team-ai-toolkit/web"
)

const (
	OrgIDKey  = "org_id"
	UserIDKey = "user_id"
)

func TenantMiddleware() web.Interceptor {
	return func(req web.Request) web.Response {
		claims := tokens.GetClaims(req)
		if claims == nil {
			return web.Err(http.StatusUnauthorized, "unauthorized", "missing claims")
		}

		audiences := claims.Audience
		if len(audiences) == 0 {
			return web.Err(http.StatusUnauthorized, "unauthorized", "missing organization")
		}

		orgID, err := uuid.Parse(audiences[0])
		if err != nil {
			return web.Err(http.StatusUnauthorized, "unauthorized", "invalid organization")
		}

		userID, err := strconv.ParseInt(claims.ID, 10, 64)
		if err != nil || userID == 0 {
			return web.Err(http.StatusUnauthorized, "unauthorized", "invalid user_id")
		}

		req.Set(OrgIDKey, orgID)
		req.Set(UserIDKey, userID)
		return web.Response{}
	}
}

func OrgID(req web.Request) uuid.UUID {
	val, exists := req.Get(OrgIDKey)
	if !exists {
		return uuid.Nil
	}
	orgID, ok := val.(uuid.UUID)
	if !ok {
		return uuid.Nil
	}
	return orgID
}

func UserID(req web.Request) int64 {
	val, exists := req.Get(UserIDKey)
	if !exists {
		return 0
	}
	id, ok := val.(int64)
	if !ok {
		return 0
	}
	return id
}
