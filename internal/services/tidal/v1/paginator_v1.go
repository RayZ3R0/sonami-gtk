package v1

import (
	"context"
	"fmt"

	"github.com/RayZ3R0/sonami-gtk/pkg/sonami"
	v1 "github.com/RayZ3R0/sonami-gtk/pkg/tidalapi/models/v1"
)

type PaginatedV1Func[Intermediary any] func(ctx context.Context, resourceID string, opts *v1.ItemsOptions) (*v1.PaginatedResponse[Intermediary], error)
type PaginatedV1Resolver[Intermediary any, Output any] func(*v1.PaginatedResponse[Intermediary]) []Output

type paginatorV1[Intermediary any, Output any] struct {
	resolver   PaginatedV1Resolver[Intermediary, Output]
	resource   PaginatedV1Func[Intermediary]
	resourceID string

	opts *v1.ItemsOptions

	items []Output
}

func NewPaginatorV1[Intermediary any, Output any](resource PaginatedV1Func[Intermediary], resourceID string, resolver PaginatedV1Resolver[Intermediary, Output], included ...string) sonami.Paginator[Output] {
	return &paginatorV1[Intermediary, Output]{
		resolver:   resolver,
		resource:   resource,
		resourceID: resourceID,
		opts: &v1.ItemsOptions{
			Limit: 20,
		},
	}
}

func (p *paginatorV1[Intermediary, Output]) IsConsumed() bool {
	return p.opts == nil
}

func (p *paginatorV1[Intermediary, Output]) NextPage() ([]Output, error) {
	if p.IsConsumed() {
		return nil, fmt.Errorf("no more pages")
	}

	paginatedResponse, err := p.resource(context.Background(), p.resourceID, p.opts)
	if err != nil {
		return nil, err
	}

	albumItems := p.resolver(paginatedResponse)

	p.items = append(p.items, albumItems...)
	if paginatedResponse.TotalNumberOfItems > len(p.items) {
		p.opts = &v1.ItemsOptions{
			Limit:  20,
			Offset: len(p.items),
		}
	} else {
		p.opts = nil
	}

	return albumItems, nil
}

func (p *paginatorV1[Intermediary, Output]) GetAll() ([]Output, error) {
	if p.IsConsumed() {
		return p.items, nil
	}

	for !p.IsConsumed() {
		_, err := p.NextPage()
		if err != nil {
			return nil, err
		}
	}

	return p.items, nil
}
