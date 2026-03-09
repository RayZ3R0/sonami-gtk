package my_collection

import (
	"github.com/RayZ3R0/sonami-gtk/internal/gettext"
	"github.com/RayZ3R0/sonami-gtk/internal/router"
	"github.com/RayZ3R0/sonami-gtk/internal/secrets"
	"github.com/RayZ3R0/sonami-gtk/internal/services/tidal/openapi"
	"github.com/RayZ3R0/sonami-gtk/internal/ui/components"
	"github.com/RayZ3R0/sonami-gtk/internal/ui/components/media_card"
	"github.com/RayZ3R0/sonami-gtk/internal/ui/pages"
	"github.com/RayZ3R0/sonami-gtk/pkg/schwifty"
	"github.com/RayZ3R0/sonami-gtk/pkg/tidalapi"
	modelopenapi "github.com/RayZ3R0/sonami-gtk/pkg/tidalapi/models/openapi"
	"github.com/RayZ3R0/sonami-gtk/pkg/tidalapi/pagination"
	"github.com/infinytum/injector"
)

func Playlists() *router.Response {
	tidal := injector.MustInject[*tidalapi.TidalAPI]()
	userId := secrets.UserID()
	if userId == "" {
		return &router.Response{
			PageTitle: gettext.Get("My Collection"),
			View:      components.AuthRequired(gettext.Get("Please sign in to view your collection")),
		}
	}

	paginator := pagination.NewPaginator(tidal.OpenAPI.V2.UserCollections.Playlists, userId, func(r *modelopenapi.Response[[]modelopenapi.Relationship]) []modelopenapi.Playlist {
		return r.Included.Playlists(r.Data...)
	}, "playlists.coverArt", "playlists.ownerProfiles")

	page, err := pages.NewPaginatedMediaCardPage(paginator, func(playlist modelopenapi.Playlist) schwifty.BaseWidgetable {
		return media_card.NewPlaylist(openapi.NewPlaylist(playlist))
	})

	return &router.Response{
		PageTitle: gettext.Get("My Playlists"),
		Error:     err,
		View:      page,
	}
}
