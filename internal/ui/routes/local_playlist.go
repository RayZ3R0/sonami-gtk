package routes

import (
	"log/slog"

	"codeberg.org/puregotk/puregotk/v4/gtk"
	"github.com/RayZ3R0/sonami-gtk/internal/cache"
	"github.com/RayZ3R0/sonami-gtk/internal/gettext"
	"github.com/RayZ3R0/sonami-gtk/internal/localdb"
	"github.com/RayZ3R0/sonami-gtk/internal/notifications"
	"github.com/RayZ3R0/sonami-gtk/internal/player"
	"github.com/RayZ3R0/sonami-gtk/internal/router"
	"github.com/RayZ3R0/sonami-gtk/internal/ui/components"
	"github.com/RayZ3R0/sonami-gtk/internal/ui/components/tracklist"
	"github.com/RayZ3R0/sonami-gtk/internal/ui/components/tracklist_header"
	"github.com/RayZ3R0/sonami-gtk/internal/ui/pages"
	"github.com/RayZ3R0/sonami-gtk/pkg/schwifty"
	. "github.com/RayZ3R0/sonami-gtk/pkg/schwifty/syntax"
	"github.com/RayZ3R0/sonami-gtk/pkg/sonami"
	"github.com/infinytum/injector"
)

var localPlaylistLogger = slog.With("module", "ui/routes", "route", "local-playlist")

func init() {
	router.Register("local-playlist/:id", LocalPlaylist)
}

// localPlaylistSource implements tracklist_header.shareablePlaybackSource.
type localPlaylistSource struct {
	pl localdb.LocalPlaylist
}

func (s localPlaylistSource) Cover(int) string              { return s.pl.CoverURL }
func (s localPlaylistSource) Route() string                 { return "local-playlist/" + s.pl.ID }
func (s localPlaylistSource) Title() string                 { return s.pl.Name }
func (s localPlaylistSource) SourceType() sonami.SourceType { return sonami.SourceTypePlaylist }
func (s localPlaylistSource) URL() string                   { return "" }

func LocalPlaylist(playlistID string) *router.Response {
	pl, err := localdb.GetPlaylist(playlistID)
	if err != nil {
		return router.FromError(gettext.Get("Playlist"), err)
	}

	trackIDs, err := localdb.GetPlaylistTrackIDs(playlistID)
	if err != nil {
		return router.FromError(pl.Name, err)
	}

	if len(trackIDs) == 0 {
		return &router.Response{
			PageTitle: pl.Name,
			View: VStack(
				components.MainContent(
					tracklist_header.NewLocalPlaylist(pl.ID, pl.Name, pl.CoverURL, 0, func() {}, func() {}).HMargin(40),
				),
				StatusPage().
					IconName("music-queue-symbolic").
					Title(gettext.Get("Empty Playlist")).
					Description(gettext.Get("Add tracks using the … menu on any track")).
					VExpand(true),
			).VMargin(20).Spacing(20),
		}
	}

	source := localPlaylistSource{*pl}

	cachedService, err := injector.Inject[*cache.CachedService]()
	if err != nil {
		return router.FromError(pl.Name, err)
	}

	tracks := cachedService.GetTrackBatch(trackIDs)

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
		return router.FromError(pl.Name, err)
	}

	playAll := func() {
		go func() {
			if err := player.PlayTracklist(source, tracks, false, 0); err != nil {
				notifications.OnToast.Notify(gettext.Get("An error occurred while playing the playlist"))
				localPlaylistLogger.Error("error playing local playlist", "error", err)
			}
		}()
	}

	shuffleAll := func() {
		go func() {
			if err := player.PlayTracklist(source, tracks, true, 0); err != nil {
				notifications.OnToast.Notify(gettext.Get("An error occurred while shuffling the playlist"))
				localPlaylistLogger.Error("error shuffling local playlist", "error", err)
			}
		}()
	}

	return &router.Response{
		PageTitle: pl.Name,
		View: VStack(
			components.MainContent(
				tracklist_header.NewLocalPlaylist(pl.ID, pl.Name, pl.CoverURL, pl.TrackCount, playAll, shuffleAll).HMargin(40),
			),
			page.VExpand(true),
		).VMargin(20).Spacing(20),
	}
}
