package auth_test

import (
	"errors"
	"fmt"
	"testing"

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
		if err == nil {
			t.Errorf("%s: expected validation error, got nil", tt.name)
			continue
		}
		if !errors.Is(err, providers.ErrValidation) {
			t.Errorf("%s: expected ErrValidation, got: %v", tt.name, err)
		}
	}
}
