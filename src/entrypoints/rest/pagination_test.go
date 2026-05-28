package rest_test

import (
	"testing"

	"github.com/educabot/team-ai-toolkit/web"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/educabot/alizia-inclusion-be/src/core/providers"
	"github.com/educabot/alizia-inclusion-be/src/entrypoints/rest"
)

func TestPage_WrapsItemsAndMoreFlag(t *testing.T) {
	items := []string{"a", "b", "c"}

	result := rest.Page(items, true)

	assert.Len(t, result.Items, 3)
	assert.True(t, result.More)
}

func TestPage_NilItemsBecomesEmptySlice(t *testing.T) {
	result := rest.Page[string](nil, false)

	assert.NotNil(t, result.Items)
	assert.Empty(t, result.Items)
	assert.False(t, result.More)
}

func TestPage_EmptySliceStaysEmpty(t *testing.T) {
	result := rest.Page([]int{}, false)

	assert.Empty(t, result.Items)
}

func TestParsePagination(t *testing.T) {
	tests := []struct {
		name       string
		limit      string
		offset     string
		wantLimit  int
		wantOffset int
		wantErr    bool
		wantErrIs  error
	}{
		{
			name:       "parses valid limit and offset",
			limit:      "25",
			offset:     "10",
			wantLimit:  25,
			wantOffset: 10,
		},
		{
			name:       "defaults to zero when no params",
			wantLimit:  0,
			wantOffset: 0,
		},
		{
			name:      "rejects invalid limit",
			limit:     "abc",
			wantErr:   true,
			wantErrIs: providers.ErrValidation,
		},
		{
			name:      "rejects invalid offset",
			limit:     "10",
			offset:    "xyz",
			wantErr:   true,
			wantErrIs: providers.ErrValidation,
		},
		{
			name:       "parses only limit",
			limit:      "50",
			wantLimit:  50,
			wantOffset: 0,
		},
	}

	for _, tt := range tests {
		req := web.NewMockRequest()
		if tt.limit != "" {
			req.Queries["limit"] = tt.limit
		}
		if tt.offset != "" {
			req.Queries["offset"] = tt.offset
		}

		p, err := rest.ParsePagination(req)

		if tt.wantErr {
			require.Error(t, err, tt.name)
			if tt.wantErrIs != nil {
				assert.ErrorIs(t, err, tt.wantErrIs, tt.name)
			}
			continue
		}
		require.NoError(t, err, tt.name)
		assert.Equal(t, tt.wantLimit, p.Limit, tt.name)
		assert.Equal(t, tt.wantOffset, p.Offset, tt.name)
	}
}
