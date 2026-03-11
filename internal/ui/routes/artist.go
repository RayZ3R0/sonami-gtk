package routes

import (
	"context"
	"log/slog"
	"strings"

	"codeberg.org/puregotk/puregotk/v4/gtk"
	"github.com/RayZ3R0/sonami-gtk/internal/gettext"
	"github.com/RayZ3R0/sonami-gtk/internal/notifications"
	"github.com/RayZ3R0/sonami-gtk/internal/player"
	"github.com/RayZ3R0/sonami-gtk/internal/router"
	"github.com/RayZ3R0/sonami-gtk/internal/ui/components"
	"github.com/RayZ3R0/sonami-gtk/internal/ui/components/tracklist_header"
	. "github.com/RayZ3R0/sonami-gtk/pkg/schwifty/syntax"
	"github.com/RayZ3R0/sonami-gtk/pkg/sonami"
	"github.com/RayZ3R0/sonami-gtk/pkg/tidalapi"
	"github.com/infinytum/injector"
)

func init() {
	router.Register("artist/:id", Artist)
}

var artistLogger = slog.With("module", "ui/routes", "route", "artist")

func Artist(artistId string) *router.Response {
	service, err := injector.Inject[sonami.Service]()
	if err != nil {
		return router.FromError(gettext.Get("Artist"), err)
	}

	artist, err := service.GetArtist(artistId)
	if err != nil {
		return router.FromError(gettext.Get("Artist"), err)
	}

	tidal := injector.MustInject[*tidalapi.TidalAPI]()
	artistPage, err := tidal.V2.Artist.Artist(context.Background(), artistId)
	if err != nil {
		return router.FromError(gettext.Get("Artist"), err)
	}

	body := VStack().Spacing(25).VMargin(20)
	for _, item := range artistPage.Items {
		if strings.Contains(strings.ToLower(item.Title), "video") {
			continue
		}
		body = body.Append(components.ForPageItem(item))
	}

	return &router.Response{
		PageTitle: gettext.Get("Artist"),
		View: VStack(
			components.MainContent(
				tracklist_header.NewArtist(artist, func() {
					go func() {
						if err := player.PlayArtistTopSongs(artistId, false, 0); err != nil {
							notifications.OnToast.Notify(gettext.Get("An error occurred while playing the top tracks"))
							albumLogger.Error("An error occurred while playing the top tracks", "error", err.Error())
						}
					}()
				}, func() {
					go func() {
						if err := player.PlayArtistTopSongs(artistId, true, 0); err != nil {
							notifications.OnToast.Notify(gettext.Get("An error occurred while playing the top tracks"))
							albumLogger.Error("An error occurred while playing the top tracks", "error", err.Error())
						}
					}()
				}).HMargin(40),
			),
			ScrolledWindow().
				Child(
					components.MainContent(
						body,
					),
				).
				Policy(gtk.PolicyNeverValue, gtk.PolicyAutomaticValue).
				VExpand(true).
				MarginTop(20),
		).VMargin(20),
	}
}
