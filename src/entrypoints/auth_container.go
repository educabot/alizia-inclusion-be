package entrypoints

import (
	authuc "github.com/educabot/alizia-inclusion-be/src/core/usecases/auth"
)

type AuthContainer struct {
	GetMe authuc.GetMe
}
