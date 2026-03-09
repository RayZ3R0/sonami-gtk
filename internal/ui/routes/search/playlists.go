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
	router.Register("search/:query/playlists", playlists)
}

func playlists(query string) *router.Response {
	tidal := injector.MustInject[*tidalapi.TidalAPI]()

	paginator := pagination.NewPaginator(tidal.OpenAPI.V2.SearchResults.Playlists, query, func(r *modelopenapi.Response[[]modelopenapi.Relationship]) []modelopenapi.Playlist {
		return r.Included.Playlists(r.Data...)
	}, "playlists.coverArt")

	page, err := pages.NewPaginatedMediaCardPage(paginator, func(playlist modelopenapi.Playlist) schwifty.BaseWidgetable {
		return media_card.NewPlaylist(openapi.NewPlaylist(playlist))
	})

	return &router.Response{
		PageTitle: gettext.Get("Search"),
		Error:     err,
		View:      page,
	}
}
