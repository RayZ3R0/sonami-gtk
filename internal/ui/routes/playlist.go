package routes

import (
	"context"
	"errors"
	"fmt"
	"time"

	"codeberg.org/dergs/tidalwave/internal/router"
	. "codeberg.org/dergs/tidalwave/pkg/gui"
	"codeberg.org/dergs/tidalwave/pkg/tidalapi"
	"github.com/diamondburned/gotk4/pkg/gtk/v4"
	"github.com/diamondburned/gotkit/gtkutil/imgutil"
	"github.com/infinytum/injector"
)

func init() {
	router.Register("playlist", Playlist)
}

func Playlist(params router.Params) *router.Response {
	playlistUUID, ok := params["uuid"].(string)
	if !ok {
		return router.FromError("Playlist", errors.New("invalid playlist ID"))
	}

	tidal := injector.MustInject[*tidalapi.TidalAPI]()

	playlist, err := tidal.V1.Playlists.Playlist(context.Background(), playlistUUID)
	if err != nil {
		return router.FromError("Playlist", err)
	}

	creator := "TIDAL"
	if playlist.Creator.Name != "" {
		creator = playlist.Creator.Name
	}

	duration := time.Duration(playlist.Duration) * time.Second

	cover := gtk.NewImage()
	cover.SetPixelSize(146)
	imgutil.AsyncGET(injector.MustInject[context.Context](), tidalapi.ImageURL(playlist.SquareImage), imgutil.ImageSetterFromImage(cover))

	return &router.Response{
		PageTitle: playlist.Title,
		View: VStack(
			HStack(
				AspectFrame(cover).CornerRadius(10).Overflow(gtk.OverflowHidden),
				VStack(
					Text(playlist.Title).
						FontSize(18).
						FontWeight(700).
						HAlign(gtk.AlignStart),
					Text(creator).
						FontSize(16).
						FontWeight(500).
						HAlign(gtk.AlignStart),
					Text(playlist.Created.Format("2006")).
						FontSize(16).
						FontWeight(500).
						HAlign(gtk.AlignStart),
					Text(fmt.Sprintf("%d Tracks (%s)", playlist.NumberOfTracks, duration.String())).
						FontSize(14).
						FontWeight(600).
						HAlign(gtk.AlignStart).
						MarginTop(10),
				).MarginLeft(20).VAlign(gtk.AlignCenter),
			),
			Spacer(),
		).HMargin(40).VMargin(20),
	}
}
