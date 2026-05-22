package auth

import (
	"fmt"

	"github.com/educabot/alizia-inclusion-be/src/core/providers"
)

var (
	errOrgIDRequired  = fmt.Errorf("%w: organization_id is required", providers.ErrValidation)
	errUserIDRequired = fmt.Errorf("%w: user_id is required", providers.ErrValidation)
)
