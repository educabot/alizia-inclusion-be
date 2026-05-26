package entrypoints_test

import (
	"fmt"

	"github.com/educabot/alizia-inclusion-be/src/core/providers"
)

var (
	errNotFound   = fmt.Errorf("%w: resource 99", providers.ErrNotFound)
	errBadRequest = fmt.Errorf("%w: invalid input", providers.ErrValidation)
)
