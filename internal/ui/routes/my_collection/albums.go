package my_collection

import (
"github.com/RayZ3R0/sonami-gtk/internal/gettext"
"github.com/RayZ3R0/sonami-gtk/internal/router"
appState "github.com/RayZ3R0/sonami-gtk/internal/state"
"github.com/RayZ3R0/sonami-gtk/internal/ui/components/media_card"
"github.com/RayZ3R0/sonami-gtk/internal/ui/pages"
. "github.com/RayZ3R0/sonami-gtk/pkg/schwifty/syntax"
"github.com/RayZ3R0/sonami-gtk/pkg/schwifty"
"github.com/RayZ3R0/sonami-gtk/pkg/sonami"
"github.com/infinytum/injector"
)

func Albums() *router.Response {
	ids, err := appState.AlbumsCache.Get()
	if err != nil {
		return router.FromError(gettext.Get("My Albums"), err)
	}

	service, err := injector.Inject[sonami.Service]()
	if err != nil {
		return router.FromError(gettext.Get("My Albums"), err)
	}

	albums := fetchAll(*ids, service.GetAlbum)

	return &router.Response{
		PageTitle: gettext.Get("My Albums"),
		View: pages.NewStaticMediaCardPage(
albums,
StatusPage().
				IconName("media-optical-cd-audio-symbolic").
				Title(gettext.Get("No Albums")).
				Description(gettext.Get("Tap the heart on an album to save it here")),
			func(a sonami.Album) schwifty.BaseWidgetable { return media_card.NewAlbum(a) },
		),
	}
}
