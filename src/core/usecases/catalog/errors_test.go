package catalog_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/educabot/alizia-inclusion-be/src/core/providers"
	"github.com/educabot/alizia-inclusion-be/src/core/usecases/catalog"
)

var (
	errDB           = errors.New("db error")
	errRampNotFound = fmt.Errorf("%w: ramp 99", providers.ErrNotFound)
	errDevNotFound  = fmt.Errorf("%w: device 99", providers.ErrNotFound)
)

func TestValidationErrors(t *testing.T) {
	tests := []struct {
		name string
		req  interface{ Validate() error }
	}{
		{"ListRamps_empty", catalog.ListRampsRequest{}},
		{"GetRamp_empty", catalog.GetRampRequest{}},
		{"ListDevices_empty", catalog.ListDevicesRequest{}},
		{"GetDevice_empty", catalog.GetDeviceRequest{}},
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
