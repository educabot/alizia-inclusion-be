package entrypoints_test

import (
	"errors"
	"fmt"

	"github.com/educabot/alizia-inclusion-be/src/core/providers"
)

var (
	errDB           = errors.New("db connection lost")
	errNotFound     = fmt.Errorf("%w: resource 99", providers.ErrNotFound)
	errBadRequest   = fmt.Errorf("%w: invalid input", providers.ErrValidation)
)
