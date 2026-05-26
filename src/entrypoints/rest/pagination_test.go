package rest_test

import (
	"errors"
	"testing"

	"github.com/educabot/team-ai-toolkit/web"

	"github.com/educabot/alizia-inclusion-be/src/core/providers"
	"github.com/educabot/alizia-inclusion-be/src/entrypoints/rest"
)

func TestPage(t *testing.T) {
	t.Run("wraps items and more flag", func(t *testing.T) {
		items := []string{"a", "b", "c"}
		result := rest.Page(items, true)
		if len(result.Items) != 3 {
			t.Errorf("expected 3 items, got %d", len(result.Items))
		}
		if !result.More {
			t.Error("expected more=true")
		}
	})

	t.Run("nil items becomes empty slice", func(t *testing.T) {
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
	})

	t.Run("empty slice stays empty", func(t *testing.T) {
		result := rest.Page([]int{}, false)
		if len(result.Items) != 0 {
			t.Errorf("expected 0 items, got %d", len(result.Items))
		}
	})
}

func TestParsePagination(t *testing.T) {
	t.Run("parses valid limit and offset", func(t *testing.T) {
		req := web.NewMockRequest()
		req.Queries["limit"] = "25"
		req.Queries["offset"] = "10"

		p, err := rest.ParsePagination(req)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if p.Limit != 25 {
			t.Errorf("expected limit 25, got %d", p.Limit)
		}
		if p.Offset != 10 {
			t.Errorf("expected offset 10, got %d", p.Offset)
		}
	})

	t.Run("defaults to zero when no params", func(t *testing.T) {
		req := web.NewMockRequest()

		p, err := rest.ParsePagination(req)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if p.Limit != 0 {
			t.Errorf("expected limit 0, got %d", p.Limit)
		}
		if p.Offset != 0 {
			t.Errorf("expected offset 0, got %d", p.Offset)
		}
	})

	t.Run("rejects invalid limit", func(t *testing.T) {
		req := web.NewMockRequest()
		req.Queries["limit"] = "abc"

		_, err := rest.ParsePagination(req)
		if err == nil {
			t.Fatal("expected error for invalid limit")
		}
		if !errors.Is(err, providers.ErrValidation) {
			t.Errorf("expected ErrValidation, got: %v", err)
		}
	})

	t.Run("rejects invalid offset", func(t *testing.T) {
		req := web.NewMockRequest()
		req.Queries["limit"] = "10"
		req.Queries["offset"] = "xyz"

		_, err := rest.ParsePagination(req)
		if err == nil {
			t.Fatal("expected error for invalid offset")
		}
		if !errors.Is(err, providers.ErrValidation) {
			t.Errorf("expected ErrValidation, got: %v", err)
		}
	})

	t.Run("parses only limit", func(t *testing.T) {
		req := web.NewMockRequest()
		req.Queries["limit"] = "50"

		p, err := rest.ParsePagination(req)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if p.Limit != 50 {
			t.Errorf("expected limit 50, got %d", p.Limit)
		}
		if p.Offset != 0 {
			t.Errorf("expected offset 0, got %d", p.Offset)
		}
	})
}
