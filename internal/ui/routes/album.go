package routes

import (
	"log/slog"

	"codeberg.org/puregotk/puregotk/v4/gtk"
	"github.com/RayZ3R0/sonami-gtk/internal/gettext"
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

var albumLogger = slog.With("module", "ui/routes", "route", "album")

func init() {
	router.Register("album/:id", Album)
}

func Album(albumId string) *router.Response {
	service, err := injector.Inject[sonami.Service]()
	if err != nil {
		return router.FromError(gettext.Get("Album"), err)
	}

	album, err := service.GetAlbum(albumId)
	if err != nil {
		return router.FromError(gettext.Get("Album"), err)
	}

	trackPaginator, err := service.GetAlbumTracks(albumId)
	if err != nil {
		return router.FromError(gettext.Get("Album"), err)
	}

	page, err := pages.NewPaginatedTracklistPage(trackPaginator, func(tl *tracklist.TrackList) schwifty.BaseWidgetable {
		tl.SetClickHandler(func(track sonami.Track, position int) {
			go func() {
				if err := player.PlayAlbum(albumId, false, position); err != nil {
					notifications.OnToast.Notify(gettext.Get("An error occurred while playing the track"))
					albumLogger.Error("An error occurred while playing the album", "error", err.Error())
				}
			}()
		})
		return tl.HMargin(30).VAlign(gtk.AlignStartValue)
	}, tracklist.PositionColumn, tracklist.TitleColumn, tracklist.ArtistsColumn, tracklist.DurationColumn, tracklist.ControlsColumn)

	if err != nil {
		return router.FromError(album.Title(), err)
	}

	return &router.Response{
		PageTitle: album.Title(),
		Error:     err,
		View: VStack(
			components.MainContent(
				tracklist_header.NewAlbum(album, func() {
					go func() {
						if err := player.PlayAlbum(albumId, false, 0); err != nil {
							notifications.OnToast.Notify(gettext.Get("An error occurred while playing the album"))
							albumLogger.Error("An error occurred while playing the album", "error", err.Error())
						}
					}()
				}, func() {
					go func() {
						if err := player.PlayAlbum(albumId, true, 0); err != nil {
							notifications.OnToast.Notify(gettext.Get("An error occurred while playing the album"))
							albumLogger.Error("An error occurred while playing the album", "error", err.Error())
						}
					}()
				}).HMargin(40),
			),
			page.VExpand(true).MarginTop(20),
		).VMargin(20),
	}
}
