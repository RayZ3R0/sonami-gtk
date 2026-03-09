package my_collection

import (
	"log/slog"

	"github.com/RayZ3R0/sonami-gtk/internal/gettext"
	"github.com/RayZ3R0/sonami-gtk/internal/notifications"
	"github.com/RayZ3R0/sonami-gtk/internal/player"
	"github.com/RayZ3R0/sonami-gtk/internal/router"
	"github.com/RayZ3R0/sonami-gtk/internal/secrets"
	"github.com/RayZ3R0/sonami-gtk/internal/services/tidal/openapi"
	"github.com/RayZ3R0/sonami-gtk/internal/ui/components"
	"github.com/RayZ3R0/sonami-gtk/internal/ui/components/tracklist"
	"github.com/RayZ3R0/sonami-gtk/internal/ui/components/tracklist_header"
	"github.com/RayZ3R0/sonami-gtk/internal/ui/pages"
	"github.com/RayZ3R0/sonami-gtk/pkg/schwifty"
	. "github.com/RayZ3R0/sonami-gtk/pkg/schwifty/syntax"
	"github.com/RayZ3R0/sonami-gtk/pkg/tidalapi"
	modelopenapi "github.com/RayZ3R0/sonami-gtk/pkg/tidalapi/models/openapi"
	"github.com/RayZ3R0/sonami-gtk/pkg/sonami"
	"codeberg.org/puregotk/puregotk/v4/gio"
	"codeberg.org/puregotk/puregotk/v4/glib"
	"codeberg.org/puregotk/puregotk/v4/gtk"
	"github.com/infinytum/injector"
)

func Tracks() *router.Response {
	logger := slog.With("module", "ui").WithGroup("ui").With("route", "my_collection").WithGroup("my_collection")
	userId := secrets.UserID()
	if userId == "" {
		return &router.Response{
			PageTitle: gettext.Get("My Collection"),
			View:      components.AuthRequired(gettext.Get("Please sign in to view your collection")),
		}
	}

	tidal := injector.MustInject[*tidalapi.TidalAPI]()
	paginator := openapi.NewPaginator(tidal.OpenAPI.V2.UserCollections.Tracks, userId, func(r *modelopenapi.Response[[]modelopenapi.Relationship]) []sonami.Track {
		results := r.Included.Tracks(r.Data...)
		tracks := make([]sonami.Track, len(results))
		for i, track := range results {
			tracks[i] = openapi.NewTrack(track)
		}
		return tracks
	}, "tracks.artists", "tracks.albums.coverArt")

	page, err := pages.NewPaginatedTracklistPage(paginator, func(tl *tracklist.TrackList) schwifty.BaseWidgetable {
		return tl.HMargin(40).VAlign(gtk.AlignStartValue)
	}, tracklist.CoverColumn, tracklist.TitleAlbumColumn, tracklist.ArtistsColumn, tracklist.DurationColumn, tracklist.ControlsColumn)

	playControlsMenu := gio.NewMenu()
	queueAllItem := gio.NewMenuItem(gettext.Get("Add My Tracks to Queue"), "win.player.queue")
	queueAllItem.SetActionAndTargetValue("win.player.queue", glib.NewVariantString("my_collection/tracks"))
	playControlsMenu.AppendItem(queueAllItem)

	return &router.Response{
		PageTitle: gettext.Get("My Tracks"),
		Error:     err,
		View: VStack(
			components.MainContent(
				tracklist_header.NewCollection(&openapi.MyTracksInfo{}, func() {
					go func() {
						tracks, err := paginator.GetAll()
						if err != nil {
							notifications.OnToast.Notify(gettext.Get("An error occurred while playing the tracks"))
							logger.Error("An error occurred while fetching the tracks", "error", err.Error())
							return
						}

						if err := player.PlayTracklist(new(openapi.MyTracksInfo), tracks, false, 0); err != nil {
							notifications.OnToast.Notify(gettext.Get("An error occurred while playing the tracks"))
							logger.Error("An error occurred while playing the tracks", "error", err.Error())
						}
					}()
				}, func() {
					go func() {
						tracks, err := paginator.GetAll()
						if err != nil {
							notifications.OnToast.Notify(gettext.Get("An error occurred while shuffling the tracks"))
							logger.Error("An error occurred while fetching the tracks", "error", err.Error())
							return
						}

						if err := player.PlayTracklist(new(openapi.MyTracksInfo), tracks, true, 0); err != nil {
							notifications.OnToast.Notify(gettext.Get("An error occurred while shuffling the tracks"))
							logger.Error("An error occurred while playing the tracks", "error", err.Error())
						}
					}()
				}).HMargin(40),
			),
			page.
				VExpand(true),
		).
			VMargin(20).
			Spacing(20),
	}
}
