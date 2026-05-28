package dashboard_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

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
		assert.ErrorIs(t, err, providers.ErrValidation, tt.name)
	}
}
