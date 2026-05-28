package providers_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/educabot/alizia-inclusion-be/src/core/providers"
)

func TestPagination_Normalize(t *testing.T) {
	tests := []struct {
		name           string
		input          providers.Pagination
		expectedLimit  int
		expectedOffset int
	}{
		{
			name:           "zero limit defaults to 50",
			input:          providers.Pagination{Limit: 0, Offset: 0},
			expectedLimit:  providers.DefaultPageLimit,
			expectedOffset: 0,
		},
		{
			name:           "negative limit defaults to 50",
			input:          providers.Pagination{Limit: -10, Offset: 0},
			expectedLimit:  providers.DefaultPageLimit,
			expectedOffset: 0,
		},
		{
			name:           "limit over max clamped to 200",
			input:          providers.Pagination{Limit: 500, Offset: 0},
			expectedLimit:  providers.MaxPageLimit,
			expectedOffset: 0,
		},
		{
			name:           "valid limit unchanged",
			input:          providers.Pagination{Limit: 25, Offset: 10},
			expectedLimit:  25,
			expectedOffset: 10,
		},
		{
			name:           "negative offset clamped to 0",
			input:          providers.Pagination{Limit: 50, Offset: -5},
			expectedLimit:  50,
			expectedOffset: 0,
		},
		{
			name:           "max boundary limit unchanged",
			input:          providers.Pagination{Limit: providers.MaxPageLimit, Offset: 0},
			expectedLimit:  providers.MaxPageLimit,
			expectedOffset: 0,
		},
		{
			name:           "limit 1 is valid",
			input:          providers.Pagination{Limit: 1, Offset: 0},
			expectedLimit:  1,
			expectedOffset: 0,
		},
	}

	for _, tc := range tests {
		got := tc.input.Normalize()
		assert.Equal(t, tc.expectedLimit, got.Limit, "%s: limit", tc.name)
		assert.Equal(t, tc.expectedOffset, got.Offset, "%s: offset", tc.name)
	}
}
