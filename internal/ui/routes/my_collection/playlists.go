package my_collection

import (
"codeberg.org/puregotk/puregotk/v4/adw"
"codeberg.org/puregotk/puregotk/v4/gtk"
"github.com/RayZ3R0/sonami-gtk/internal/gettext"
"github.com/RayZ3R0/sonami-gtk/internal/localdb"
"github.com/RayZ3R0/sonami-gtk/internal/router"
appState "github.com/RayZ3R0/sonami-gtk/internal/state"
"github.com/RayZ3R0/sonami-gtk/internal/ui/components"
"github.com/RayZ3R0/sonami-gtk/internal/ui/components/media_card"
. "github.com/RayZ3R0/sonami-gtk/pkg/schwifty/syntax"
"github.com/RayZ3R0/sonami-gtk/pkg/sonami"
"github.com/infinytum/injector"
)

func Playlists() *router.Response {
	localPlaylists, _ := localdb.GetAllPlaylists()

	var tidalPlaylists []sonami.Playlist
	if ids, err := appState.PlaylistsCache.Get(); err == nil && len(*ids) > 0 {
		if service, err := injector.Inject[sonami.Service](); err == nil {
			tidalPlaylists = fetchAll(*ids, service.GetPlaylist)
		}
	}

	if len(localPlaylists) == 0 && len(tidalPlaylists) == 0 {
		return &router.Response{
			PageTitle: gettext.Get("My Playlists"),
			View: StatusPage().
				IconName("music-queue-symbolic").
				Title(gettext.Get("No Playlists")).
				Description(gettext.Get("Create playlists or tap the heart on a playlist to save it here")),
		}
	}

	var sections []any

	if len(localPlaylists) > 0 {
		grid := WrapBox().
			HMargin(40).
			VAlign(gtk.AlignStartValue).
			Justify(adw.JustifyFillValue).
			JustifyLastLine(true)()
		for _, p := range localPlaylists {
			grid.Append(CenterBox().CenterWidget(media_card.NewLocalPlaylist(p)).ToGTK())
		}
		sections = append(sections,
Label(gettext.Get("My Playlists")).HMargin(40).WithCSSClass("title-2").HAlign(gtk.AlignStartValue),
Widget(&grid.Widget),
		)
	}

	if len(tidalPlaylists) > 0 {
		title := gettext.Get("Saved Playlists")
		if len(localPlaylists) == 0 {
			title = gettext.Get("My Playlists")
		}
		grid := WrapBox().
			HMargin(40).
			VAlign(gtk.AlignStartValue).
			Justify(adw.JustifyFillValue).
			JustifyLastLine(true)()
		for _, p := range tidalPlaylists {
			grid.Append(CenterBox().CenterWidget(media_card.NewPlaylist(p)).ToGTK())
		}
		sections = append(sections,
Label(title).HMargin(40).WithCSSClass("title-2").HAlign(gtk.AlignStartValue),
Widget(&grid.Widget),
		)
	}

	return &router.Response{
		PageTitle: gettext.Get("My Playlists"),
		View: ScrolledWindow().
			Child(components.MainContent(VStack(sections...).VMargin(20).Spacing(24))).
			Policy(gtk.PolicyNeverValue, gtk.PolicyAutomaticValue),
	}
}
