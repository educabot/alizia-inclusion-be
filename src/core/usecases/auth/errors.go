package auth

import (
	"fmt"

	"github.com/educabot/alizia-inclusion-be/src/core/providers"
)

var (
	errEmailRequired    = fmt.Errorf("%w: email is required", providers.ErrValidation)
	errPasswordRequired = fmt.Errorf("%w: password is required", providers.ErrValidation)
	errOrgIDRequired    = fmt.Errorf("%w: organization_id is required", providers.ErrValidation)
	errUserIDRequired   = fmt.Errorf("%w: user_id is required", providers.ErrValidation)
)
