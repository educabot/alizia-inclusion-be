package rest_test

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/educabot/alizia-inclusion-be/src/core/providers"
	"github.com/educabot/alizia-inclusion-be/src/entrypoints/rest"
)

func TestHandleError(t *testing.T) {
	tests := []struct {
		name           string
		err            error
		expectedStatus int
	}{
		{
			name:           "ErrProfileNotFound returns 404",
			err:            fmt.Errorf("student 5: %w", providers.ErrProfileNotFound),
			expectedStatus: http.StatusNotFound,
		},
		{
			name:           "ErrServiceUnavailable returns 503",
			err:            fmt.Errorf("AI call: %w", providers.ErrServiceUnavailable),
			expectedStatus: http.StatusServiceUnavailable,
		},
		{
			name:           "ErrClassroomNotFound returns 404",
			err:            fmt.Errorf("classroom 3: %w", providers.ErrClassroomNotFound),
			expectedStatus: http.StatusNotFound,
		},
		{
			name:           "ErrAdaptationNotFound returns 404",
			err:            fmt.Errorf("adaptation 7: %w", providers.ErrAdaptationNotFound),
			expectedStatus: http.StatusNotFound,
		},
		{
			name:           "ErrValidation returns 400",
			err:            fmt.Errorf("bad field: %w", providers.ErrValidation),
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "ErrNotFound returns 404",
			err:            fmt.Errorf("resource: %w", providers.ErrNotFound),
			expectedStatus: http.StatusNotFound,
		},
		{
			name:           "ErrUnauthorized returns 401",
			err:            fmt.Errorf("no creds: %w", providers.ErrUnauthorized),
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "ErrForbidden returns 403",
			err:            fmt.Errorf("no access: %w", providers.ErrForbidden),
			expectedStatus: http.StatusForbidden,
		},
		{
			name:           "ErrDuplicate returns 409",
			err:            fmt.Errorf("already exists: %w", providers.ErrDuplicate),
			expectedStatus: http.StatusConflict,
		},
		{
			name:           "unknown error returns 500",
			err:            fmt.Errorf("something unexpected"),
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tc := range tests {
		resp := rest.HandleError(tc.err)
		assert.Equal(t, tc.expectedStatus, resp.Status, tc.name)
	}
}
