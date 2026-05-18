package entrypoints

import (
	"fmt"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"github.com/educabot/team-ai-toolkit/tokens"
	"github.com/educabot/team-ai-toolkit/web"

	"github.com/educabot/alizia-inclusion-be/src/core/entities"
	authuc "github.com/educabot/alizia-inclusion-be/src/core/usecases/auth"
	"github.com/educabot/alizia-inclusion-be/src/entrypoints/middleware"
	"github.com/educabot/alizia-inclusion-be/src/entrypoints/rest"
)

type loginBody struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type loginResponse struct {
	Token string       `json:"token"`
	User  userResponse `json:"user"`
}

type userResponse struct {
	ID    int64  `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Role  string `json:"role"`
}

func mapUser(u entities.User) userResponse {
	return userResponse{
		ID:    u.ID,
		Name:  u.Name,
		Email: u.Email,
		Role:  u.Role,
	}
}

func mapUsers(us []entities.User) []userResponse {
	out := make([]userResponse, len(us))
	for i, u := range us {
		out[i] = mapUser(u)
	}
	return out
}

func (c *AuthContainer) HandleLogin(req web.Request) web.Response {
	var body loginBody
	if err := req.BindJSON(&body); err != nil {
		return rest.HandleError(err)
	}

	user, err := c.LoginUC.Execute(req.Context(), authuc.LoginRequest{
		Email:    body.Email,
		Password: body.Password,
	})
	if err != nil {
		return rest.HandleError(err)
	}

	claims := tokens.Claims{
		ID:    fmt.Sprintf("%d", user.ID),
		Name:  user.Name,
		Email: user.Email,
		Roles: []string{user.Role},
		RegisteredClaims: jwt.RegisteredClaims{
			Audience: jwt.ClaimStrings{user.OrganizationID.String()},
		},
	}
	token, err := c.Toker.CreateWithClaims(claims, 24*time.Hour)
	if err != nil {
		return web.Err(http.StatusInternalServerError, "token_error", "failed to create token")
	}

	return web.OK(loginResponse{
		Token: token,
		User:  mapUser(*user),
	})
}

func (c *AuthContainer) HandleGetMe(req web.Request) web.Response {
	user, err := c.GetMe.Execute(req.Context(), authuc.GetMeRequest{
		OrgID:  middleware.OrgID(req),
		UserID: middleware.UserID(req),
	})
	if err != nil {
		return rest.HandleError(err)
	}
	return web.OK(mapUser(*user))
}
