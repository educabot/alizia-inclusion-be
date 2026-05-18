package rest

import (
	"net/http"

	bcerrors "github.com/educabot/team-ai-toolkit/errors"
	"github.com/educabot/team-ai-toolkit/web"

	"github.com/educabot/alizia-inclusion-be/src/core/providers"
)

func HandleError(err error) web.Response {
	switch {
	case bcerrors.Is(err, providers.ErrProfileNotFound):
		return web.Err(http.StatusNotFound, "profile_not_found", err.Error())
	case bcerrors.Is(err, providers.ErrServiceUnavailable):
		return web.Err(http.StatusServiceUnavailable, "service_unavailable", err.Error())
	case bcerrors.Is(err, providers.ErrInvalidCredentials):
		return web.Err(http.StatusUnauthorized, "invalid_credentials", err.Error())
	case bcerrors.Is(err, providers.ErrClassroomNotFound):
		return web.Err(http.StatusNotFound, "classroom_not_found", err.Error())
	case bcerrors.Is(err, providers.ErrAdaptationNotFound):
		return web.Err(http.StatusNotFound, "adaptation_not_found", err.Error())
	default:
		return bcerrors.HandleError(err)
	}
}
