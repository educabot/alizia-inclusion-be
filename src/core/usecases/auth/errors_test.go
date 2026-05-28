package auth_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/educabot/alizia-inclusion-be/src/core/providers"
	"github.com/educabot/alizia-inclusion-be/src/core/usecases/auth"
)

var errUserNotFound = fmt.Errorf("%w: user 99", providers.ErrNotFound)

func TestValidationErrors(t *testing.T) {
	tests := []struct {
		name string
		req  interface{ Validate() error }
	}{
		{"GetMe_empty", auth.GetMeRequest{}},
	}
	for _, tt := range tests {
		err := tt.req.Validate()
		assert.Error(t, err, tt.name)
		assert.ErrorIs(t, err, providers.ErrValidation, tt.name)
	}
}
