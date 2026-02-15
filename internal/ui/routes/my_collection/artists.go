package my_collection

import (
	"codeberg.org/dergs/tonearm/internal/gettext"
	"codeberg.org/dergs/tonearm/internal/router"
	"codeberg.org/dergs/tonearm/internal/secrets"
	"codeberg.org/dergs/tonearm/internal/ui/components"
	"codeberg.org/dergs/tonearm/internal/ui/components/media_card"
	"codeberg.org/dergs/tonearm/internal/ui/pages"
	"codeberg.org/dergs/tonearm/pkg/schwifty"
	"codeberg.org/dergs/tonearm/pkg/tidalapi"
	"codeberg.org/dergs/tonearm/pkg/tidalapi/models/openapi"
	"codeberg.org/dergs/tonearm/pkg/tidalapi/pagination"
	"github.com/infinytum/injector"
)

func Artists() *router.Response {
	tidal := injector.MustInject[*tidalapi.TidalAPI]()
	userId := secrets.UserID()
	if userId == "" {
		return &router.Response{
			PageTitle: gettext.Get("My Collection"),
			View: components.AuthRequired(gettext.Get("Please sign in to view your collection")),
		}
	}

	paginator := pagination.NewPaginator(tidal.OpenAPI.V2.UserCollections.Artists, userId, func(r *openapi.Response[[]openapi.Relationship]) []openapi.Artist {
		return r.Included.Artists(r.Data...)
	}, "artists.profileArt")

	page, err := pages.NewPaginatedMediaCardPage(paginator, func(artist openapi.Artist) schwifty.BaseWidgetable {
		return media_card.NewArtist(&artist)
	})

	return &router.Response{
		PageTitle: gettext.Get("My Artists"),
		Error:     err,
		View:      page,
	}
}
