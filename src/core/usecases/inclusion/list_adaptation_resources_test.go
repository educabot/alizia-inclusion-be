package inclusion_test

import (
	"context"
	"errors"
	"testing"

	"github.com/educabot/alizia-inclusion-be/src/core/entities"
	"github.com/educabot/alizia-inclusion-be/src/core/providers"
	"github.com/educabot/alizia-inclusion-be/src/core/providers/mocks"
	"github.com/educabot/alizia-inclusion-be/src/core/usecases/inclusion"
	"github.com/educabot/alizia-inclusion-be/src/testutil"
)

func TestListAdaptationResources(t *testing.T) {
	ctx := context.Background()

	t.Run("returns resources", func(t *testing.T) {
		expected := []entities.AdaptationResource{
			testutil.NewAdaptationResource(1, 10),
			testutil.NewAdaptationResource(2, 10),
		}
		mock := &mocks.MockAdaptationResourceProvider{
			ListByAdaptationFn: func(_ context.Context, adaptationID int64) ([]entities.AdaptationResource, error) {
				if adaptationID != 10 {
					t.Errorf("expected adaptationID 10, got %d", adaptationID)
				}
				return expected, nil
			},
		}

		got, err := inclusion.NewListAdaptationResources(mock).Execute(ctx, inclusion.ListAdaptationResourcesRequest{
			AdaptationID: 10,
		})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(got) != 2 {
			t.Errorf("got %d resources, want 2", len(got))
		}
	})

	t.Run("rejects zero adaptation_id", func(t *testing.T) {
		mock := &mocks.MockAdaptationResourceProvider{}
		_, err := inclusion.NewListAdaptationResources(mock).Execute(ctx, inclusion.ListAdaptationResourcesRequest{
			AdaptationID: 0,
		})
		if !errors.Is(err, providers.ErrValidation) {
			t.Errorf("expected ErrValidation, got: %v", err)
		}
	})
}
