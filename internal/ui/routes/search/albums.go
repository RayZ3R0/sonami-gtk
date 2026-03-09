package search

import (
	"github.com/RayZ3R0/sonami-gtk/internal/gettext"
	"github.com/RayZ3R0/sonami-gtk/internal/router"
	"github.com/RayZ3R0/sonami-gtk/internal/services/tidal/openapi"
	"github.com/RayZ3R0/sonami-gtk/internal/ui/components/media_card"
	"github.com/RayZ3R0/sonami-gtk/internal/ui/pages"
	"github.com/RayZ3R0/sonami-gtk/pkg/schwifty"
	"github.com/RayZ3R0/sonami-gtk/pkg/tidalapi"
	modelopenapi "github.com/RayZ3R0/sonami-gtk/pkg/tidalapi/models/openapi"
	"github.com/RayZ3R0/sonami-gtk/pkg/tidalapi/pagination"
	"github.com/infinytum/injector"
)

func init() {
	router.Register("search/:query/albums", albums)
}

func albums(query string) *router.Response {
	tidal := injector.MustInject[*tidalapi.TidalAPI]()

	paginator := pagination.NewPaginator(tidal.OpenAPI.V2.SearchResults.Albums, query, func(r *modelopenapi.Response[[]modelopenapi.Relationship]) []modelopenapi.Album {
		return r.Included.Albums(r.Data...)
	}, "albums.coverArt", "albums.artists")

	page, err := pages.NewPaginatedMediaCardPage(paginator, func(album modelopenapi.Album) schwifty.BaseWidgetable {
		return media_card.NewAlbum(openapi.NewAlbum(album))
	})

	return &router.Response{
		PageTitle: gettext.Get("Search"),
		Error:     err,
		View:      page,
	}
}
