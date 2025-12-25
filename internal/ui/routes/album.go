package routes

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"codeberg.org/dergs/tidalwave/internal/router"
	"codeberg.org/dergs/tidalwave/internal/ui/components/tracklist"
	. "codeberg.org/dergs/tidalwave/pkg/gui"
	"codeberg.org/dergs/tidalwave/pkg/tidalapi"
	"github.com/diamondburned/gotk4/pkg/gtk/v4"
	"github.com/diamondburned/gotkit/gtkutil/imgutil"
	"github.com/infinytum/injector"
)

func init() {
	router.Register("album", Album)
}

func Album(params router.Params) *router.Response {
	albumId, ok := params["id"].(int)
	if !ok {
		return router.FromError("Album", errors.New("invalid album ID"))
	}

	tidal := injector.MustInject[*tidalapi.TidalAPI]()

	album, err := tidal.OpenAPI.V2.Albums.Album(context.Background(), albumId, "coverArt", "artists", "coverArt")
	if err != nil {
		return router.FromError("Album", err)
	}

	// TODO: Handle pagination with scroll events
	items, err := tidal.OpenAPI.V2.Albums.Items(context.Background(), albumId, "", "items", "items.artists")
	if err != nil {
		return router.FromError("Album", err)
	}

	artists := []string{}
	for _, artist := range album.Included.PlainArtists(album.Data.Relationships.Artists.Data...) {
		artists = append(artists, artist.Attributes.Name)
	}

	cover := gtk.NewImage()
	cover.SetPixelSize(146)
	for _, artwork := range album.Included.PlainArtworks(album.Data.Relationships.CoverArt.Data...) {
		imgutil.AsyncGET(injector.MustInject[context.Context](), artwork.Attributes.Files.AtLeast(320).Href, imgutil.ImageSetterFromImage(cover))
		break
	}

	list := tracklist.NewTrackList(
		tracklist.PositionColumn,
		tracklist.TitleColumn,
		tracklist.ArtistsColumn,
		tracklist.DurationColumn,
		tracklist.BoxColumn,
		tracklist.ControlsColumn,
	)

	for _, track := range items.Included.Tracks(items.Data...) {
		list.AddTrack(&track)
	}

	scroll := gtk.NewScrolledWindow()
	scroll.SetPolicy(gtk.PolicyNever, gtk.PolicyAutomatic)
	scroll.SetChild(list.SetTitle(""))

	return &router.Response{
		PageTitle: album.Data.Attributes.Title,
		View: VStack(
			HStack(
				AspectFrame(cover).CornerRadius(10).Overflow(gtk.OverflowHidden),
				VStack(
					Text(album.Data.Attributes.Title).
						FontSize(18).
						FontWeight(700).
						HAlign(gtk.AlignStart),
					Text(strings.Join(artists, ", ")).
						FontSize(16).
						FontWeight(500).
						HAlign(gtk.AlignStart),
					Text(album.Data.Attributes.ReleaseDate.Format("2006")).
						FontSize(16).
						FontWeight(500).
						HAlign(gtk.AlignStart),
					Text(fmt.Sprintf("%d Tracks (%s)", album.Data.Attributes.NumberOfItems, tidalapi.FormatDuration(int(album.Data.Attributes.Duration.Seconds())))).
						FontSize(14).
						FontWeight(600).
						HAlign(gtk.AlignStart).
						MarginTop(10),
				).MarginLeft(20).VAlign(gtk.AlignCenter),
			),
			Wrapper(scroll).VExpand(true).MarginTop(20),
		).HMargin(40).VMargin(20),
	}
}
