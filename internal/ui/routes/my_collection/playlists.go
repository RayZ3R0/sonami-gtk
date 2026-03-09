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

func Playlists() *router.Response {
	ids, err := appState.PlaylistsCache.Get()
	if err != nil {
		return router.FromError(gettext.Get("My Playlists"), err)
	}

	service, err := injector.Inject[sonami.Service]()
	if err != nil {
		return router.FromError(gettext.Get("My Playlists"), err)
	}

	playlists := fetchAll(*ids, service.GetPlaylist)

	return &router.Response{
		PageTitle: gettext.Get("My Playlists"),
		View: pages.NewStaticMediaCardPage(
playlists,
StatusPage().
				IconName("music-queue-symbolic").
				Title(gettext.Get("No Playlists")).
				Description(gettext.Get("Tap the heart on a playlist to save it here")),
			func(p sonami.Playlist) schwifty.BaseWidgetable { return media_card.NewPlaylist(p) },
		),
	}
}
