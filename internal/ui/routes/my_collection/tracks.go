package my_collection

import (
	"log/slog"

	"codeberg.org/puregotk/puregotk/v4/gtk"
	"github.com/RayZ3R0/sonami-gtk/internal/gettext"
	"github.com/RayZ3R0/sonami-gtk/internal/notifications"
	"github.com/RayZ3R0/sonami-gtk/internal/player"
	"github.com/RayZ3R0/sonami-gtk/internal/router"
	appState "github.com/RayZ3R0/sonami-gtk/internal/state"
	"github.com/RayZ3R0/sonami-gtk/internal/ui/components"
	"github.com/RayZ3R0/sonami-gtk/internal/ui/components/tracklist"
	"github.com/RayZ3R0/sonami-gtk/internal/ui/components/tracklist_header"
	"github.com/RayZ3R0/sonami-gtk/internal/ui/pages"
	"github.com/RayZ3R0/sonami-gtk/pkg/schwifty"
	. "github.com/RayZ3R0/sonami-gtk/pkg/schwifty/syntax"
	"github.com/RayZ3R0/sonami-gtk/pkg/sonami"
	"github.com/infinytum/injector"
)

var tracksLogger = slog.With("module", "ui/routes/my_collection", "route", "tracks")

// myTracksSource satisfies tracklist_header.shareablePlaybackSource for the "Liked Tracks" header.
type myTracksSource struct{}

func (s myTracksSource) Cover(int) string              { return "" }
func (s myTracksSource) Route() string                 { return "my-collection/tracks" }
func (s myTracksSource) Title() string                 { return gettext.Get("Liked Tracks") }
func (s myTracksSource) SourceType() sonami.SourceType { return sonami.SourceTypePlaylist }
func (s myTracksSource) URL() string                   { return "" }

func Tracks() *router.Response {
	ids, err := appState.TracksCache.Get()
	if err != nil {
		return router.FromError(gettext.Get("My Tracks"), err)
	}

	if len(*ids) == 0 {
		return &router.Response{
			PageTitle: gettext.Get("My Tracks"),
			View: StatusPage().
				IconName("heart-outline-thick-symbolic").
				Title(gettext.Get("No Tracks")).
				Description(gettext.Get("Tap the heart on a track to save it here")),
		}
	}

	service, err := injector.Inject[sonami.Service]()
	if err != nil {
		return router.FromError(gettext.Get("My Tracks"), err)
	}

	tracks := fetchAll(*ids, service.GetTrack)

	paginator := sonami.NewArrayPaginator(tracks)
	page, err := pages.NewPaginatedTracklistPage(
		paginator,
		func(tl *tracklist.TrackList) schwifty.BaseWidgetable {
			return tl.HMargin(40).VAlign(gtk.AlignStartValue)
		},
		tracklist.CoverColumn,
		tracklist.TitleAlbumColumn,
		tracklist.ArtistsColumn,
		tracklist.DurationColumn,
		tracklist.ControlsColumn,
	)
	if err != nil {
		return router.FromError(gettext.Get("My Tracks"), err)
	}

	playAll := func() {
		go func() {
			if err := player.PlayTracklist(myTracksSource{}, tracks, false, 0); err != nil {
				notifications.OnToast.Notify(gettext.Get("An error occurred while playing your tracks"))
				tracksLogger.Error("error playing liked tracks", "error", err)
			}
		}()
	}

	shuffleAll := func() {
		go func() {
			if err := player.PlayTracklist(myTracksSource{}, tracks, true, 0); err != nil {
				notifications.OnToast.Notify(gettext.Get("An error occurred while shuffling your tracks"))
				tracksLogger.Error("error shuffling liked tracks", "error", err)
			}
		}()
	}

	return &router.Response{
		PageTitle: gettext.Get("My Tracks"),
		View: VStack(
			components.MainContent(
				tracklist_header.NewCollection(myTracksSource{}, playAll, shuffleAll).HMargin(40),
			),
			page.VExpand(true),
		).VMargin(20).Spacing(20),
	}
}
