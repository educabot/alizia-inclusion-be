package rest

import (
	"fmt"
	"strconv"

	"github.com/educabot/team-ai-toolkit/web"

	"github.com/educabot/alizia-inclusion-be/src/core/providers"
)

type PaginatedResponse[T any] struct {
	Items []T  `json:"items"`
	More  bool `json:"more"`
}

func Page[T any](items []T, more bool) PaginatedResponse[T] {
	if items == nil {
		items = []T{}
	}
	return PaginatedResponse[T]{Items: items, More: more}
}

func ParsePagination(req web.Request) (providers.Pagination, error) {
	var p providers.Pagination
	if v := req.Query("limit"); v != "" {
		n, err := strconv.Atoi(v)
		if err != nil {
			return p, fmt.Errorf("%w: invalid limit", providers.ErrValidation)
		}
		p.Limit = n
	}
	if v := req.Query("offset"); v != "" {
		n, err := strconv.Atoi(v)
		if err != nil {
			return p, fmt.Errorf("%w: invalid offset", providers.ErrValidation)
		}
		p.Offset = n
	}
	return p, nil
}
