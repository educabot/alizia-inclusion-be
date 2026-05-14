package catalog

import (
	"fmt"

	"github.com/educabot/alizia-inclusion-be/src/core/providers"
)

var (
	errOrgIDRequired    = fmt.Errorf("%w: organization_id is required", providers.ErrValidation)
	errRampIDRequired   = fmt.Errorf("%w: ramp_id is required", providers.ErrValidation)
	errDeviceIDRequired = fmt.Errorf("%w: device_id is required", providers.ErrValidation)
)
