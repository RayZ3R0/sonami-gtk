package my_collection

import (
	"codeberg.org/dergs/tonearm/internal/gettext"
	"codeberg.org/dergs/tonearm/internal/router"
	"codeberg.org/dergs/tonearm/internal/secrets"
	"codeberg.org/dergs/tonearm/internal/ui/components/media_card"
	"codeberg.org/dergs/tonearm/internal/ui/pages"
	"codeberg.org/dergs/tonearm/pkg/schwifty"
	. "codeberg.org/dergs/tonearm/pkg/schwifty/syntax"
	"codeberg.org/dergs/tonearm/pkg/tidalapi"
	"codeberg.org/dergs/tonearm/pkg/tidalapi/models/openapi"
	"codeberg.org/dergs/tonearm/pkg/tidalapi/pagination"
	"github.com/infinytum/injector"
)

func Albums() *router.Response {
	tidal := injector.MustInject[*tidalapi.TidalAPI]()
	userId := secrets.UserID()
	if userId == "" {
		return &router.Response{
			PageTitle: gettext.Get("My Collection"),
			View:      Label(gettext.Get("Please log in to view your collection")),
		}
	}

	paginator := pagination.NewPaginator(tidal.OpenAPI.V2.UserCollections.Albums, userId, func(r *openapi.Response[[]openapi.Relationship]) []openapi.Album {
		return r.Included.Albums(r.Data...)
	}, "albums.coverArt", "albums.artists")

	page, err := pages.NewPaginatedMediaCardPage(paginator, func(album openapi.Album) schwifty.BaseWidgetable {
		return media_card.NewAlbum(&album)
	})

	return &router.Response{
		PageTitle: gettext.Get("My Albums"),
		Error:     err,
		View:      page,
	}
}
