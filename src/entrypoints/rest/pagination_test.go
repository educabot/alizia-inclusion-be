package rest_test

import (
	"errors"
	"testing"

	"github.com/educabot/team-ai-toolkit/web"

	"github.com/educabot/alizia-inclusion-be/src/core/providers"
	"github.com/educabot/alizia-inclusion-be/src/entrypoints/rest"
)

func TestPage_WrapsItemsAndMoreFlag(t *testing.T) {
	items := []string{"a", "b", "c"}

	result := rest.Page(items, true)

	if len(result.Items) != 3 {
		t.Errorf("expected 3 items, got %d", len(result.Items))
	}
	if !result.More {
		t.Error("expected more=true")
	}
}

func TestPage_NilItemsBecomesEmptySlice(t *testing.T) {
	result := rest.Page[string](nil, false)

	if result.Items == nil {
		t.Error("expected non-nil empty slice")
	}
	if len(result.Items) != 0 {
		t.Errorf("expected 0 items, got %d", len(result.Items))
	}
	if result.More {
		t.Error("expected more=false")
	}
}

func TestPage_EmptySliceStaysEmpty(t *testing.T) {
	result := rest.Page([]int{}, false)

	if len(result.Items) != 0 {
		t.Errorf("expected 0 items, got %d", len(result.Items))
	}
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
			if err == nil {
				t.Errorf("%s: expected error, got nil", tt.name)
				continue
			}
			if tt.wantErrIs != nil && !errors.Is(err, tt.wantErrIs) {
				t.Errorf("%s: expected %v, got: %v", tt.name, tt.wantErrIs, err)
			}
			continue
		}
		if err != nil {
			t.Errorf("%s: unexpected error: %v", tt.name, err)
			continue
		}
		if p.Limit != tt.wantLimit {
			t.Errorf("%s: expected limit %d, got %d", tt.name, tt.wantLimit, p.Limit)
		}
		if p.Offset != tt.wantOffset {
			t.Errorf("%s: expected offset %d, got %d", tt.name, tt.wantOffset, p.Offset)
		}
	}
}
