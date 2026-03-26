package my_collection

import (
	"github.com/RayZ3R0/sonami-gtk/internal/cache"
	"github.com/RayZ3R0/sonami-gtk/internal/gettext"
	"github.com/RayZ3R0/sonami-gtk/internal/router"
	appState "github.com/RayZ3R0/sonami-gtk/internal/state"
	"github.com/RayZ3R0/sonami-gtk/internal/ui/components/media_card"
	"github.com/RayZ3R0/sonami-gtk/internal/ui/pages"
	"github.com/RayZ3R0/sonami-gtk/pkg/schwifty"
	. "github.com/RayZ3R0/sonami-gtk/pkg/schwifty/syntax"
	"github.com/RayZ3R0/sonami-gtk/pkg/sonami"
	"github.com/infinytum/injector"
)

func Artists() *router.Response {
	ids, err := appState.ArtistsCache.Get()
	if err != nil {
		return router.FromError(gettext.Get("My Artists"), err)
	}

	cachedService, err := injector.Inject[*cache.CachedService]()
	if err != nil {
		return router.FromError(gettext.Get("My Artists"), err)
	}

	artists := cachedService.GetArtistBatch(*ids)

	// sonami.Artist implements sonami.ArtistInfo (the parameter of media_card.NewArtist).
	return &router.Response{
		PageTitle: gettext.Get("My Artists"),
		View: pages.NewStaticMediaCardPage(
			artists,
			StatusPage().
				IconName("music-artist2-symbolic").
				Title(gettext.Get("No Artists")).
				Description(gettext.Get("Tap the heart on an artist to save them here")),
			func(a sonami.Artist) schwifty.BaseWidgetable { return media_card.NewArtist(a) },
		),
	}
}
