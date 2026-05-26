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
