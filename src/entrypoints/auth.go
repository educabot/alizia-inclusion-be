package entrypoints

import (
	"github.com/educabot/team-ai-toolkit/web"

	"github.com/educabot/alizia-inclusion-be/src/core/entities"
	authuc "github.com/educabot/alizia-inclusion-be/src/core/usecases/auth"
	"github.com/educabot/alizia-inclusion-be/src/entrypoints/middleware"
	"github.com/educabot/alizia-inclusion-be/src/entrypoints/rest"
)

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
