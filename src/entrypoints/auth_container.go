package entrypoints

import (
	authuc "github.com/educabot/alizia-inclusion-be/src/core/usecases/auth"
	"github.com/educabot/team-ai-toolkit/tokens"
)

type AuthContainer struct {
	Toker   tokens.Toker
	LoginUC authuc.Login
	GetMe   authuc.GetMe
}
