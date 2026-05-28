package dashboard_test

import (
	"errors"
	"testing"

	"github.com/educabot/alizia-inclusion-be/src/core/providers"
	"github.com/educabot/alizia-inclusion-be/src/core/usecases/dashboard"
)

func TestValidationErrors(t *testing.T) {
	tests := []struct {
		name string
		req  interface{ Validate() error }
	}{
		{"GetMetrics_empty", dashboard.GetMetricsRequest{}},
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
