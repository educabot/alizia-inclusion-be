package inclusion

import (
	"context"

	"github.com/educabot/alizia-inclusion-be/src/core/entities"
	"github.com/educabot/alizia-inclusion-be/src/core/providers"
)

type ListAdaptationResourcesRequest struct {
	AdaptationID int64
}

func (r ListAdaptationResourcesRequest) Validate() error {
	if r.AdaptationID <= 0 {
		return errAdaptationIDRequired
	}
	return nil
}

type ListAdaptationResources interface {
	Execute(ctx context.Context, req ListAdaptationResourcesRequest) ([]entities.AdaptationResource, error)
}

type listAdaptationResourcesImpl struct {
	resources providers.AdaptationResourceProvider
}

func NewListAdaptationResources(resources providers.AdaptationResourceProvider) ListAdaptationResources {
	return &listAdaptationResourcesImpl{resources: resources}
}

func (uc *listAdaptationResourcesImpl) Execute(ctx context.Context, req ListAdaptationResourcesRequest) ([]entities.AdaptationResource, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	return uc.resources.ListByAdaptation(ctx, req.AdaptationID)
}
