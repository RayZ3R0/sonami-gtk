package routes

import (
	"context"
	"errors"
	"fmt"

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

	playlist, err := tidal.OpenAPI.V2.Playlists.Playlist(context.Background(), playlistUUID, "coverArt", "ownerProfiles")
	if err != nil {
		return router.FromError("Playlist", err)
	}

	creator := "TIDAL"
	cover := gtk.NewImage()
	cover.SetPixelSize(146)

	for _, item := range playlist.Included {
		if item.Attributes.Artworks != nil {
			imgutil.AsyncGET(injector.MustInject[context.Context](), item.Attributes.Artworks.Files.AtLeast(320).Href, imgutil.ImageSetterFromImage(cover))
		}
		if item.Attributes.Artist != nil {
			creator = item.Attributes.Artist.Name
		}
	}

	return &router.Response{
		PageTitle: playlist.Data.Attributes.Name,
		View: VStack(
			HStack(
				AspectFrame(cover).CornerRadius(10).Overflow(gtk.OverflowHidden),
				VStack(
					Text(playlist.Data.Attributes.Name).
						FontSize(18).
						FontWeight(700).
						HAlign(gtk.AlignStart),
					Text(creator).
						FontSize(16).
						FontWeight(500).
						HAlign(gtk.AlignStart),
					Text(playlist.Data.Attributes.CreatedAt.Format("2006")).
						FontSize(16).
						FontWeight(500).
						HAlign(gtk.AlignStart),
					Text(fmt.Sprintf("%d Tracks (%s)", playlist.Data.Attributes.NumberOfItems, tidalapi.FormatDuration(int(playlist.Data.Attributes.Duration.Seconds())))).
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
