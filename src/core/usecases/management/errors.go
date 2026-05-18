package management

import (
	"fmt"

	"github.com/educabot/alizia-inclusion-be/src/core/providers"
)

var (
	errOrgIDRequired       = fmt.Errorf("%w: organization_id is required", providers.ErrValidation)
	errClassroomIDRequired = fmt.Errorf("%w: classroom_id is required", providers.ErrValidation)
	errNameRequired        = fmt.Errorf("%w: name is required", providers.ErrValidation)
)
