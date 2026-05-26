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
		t.Run(tt.name, func(t *testing.T) {
			err := tt.req.Validate()
			if err == nil {
				t.Fatal("expected validation error, got nil")
			}
			if !errors.Is(err, providers.ErrValidation) {
				t.Errorf("expected ErrValidation, got: %v", err)
			}
		})
	}
}
