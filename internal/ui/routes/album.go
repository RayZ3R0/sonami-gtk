package routes

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"strings"

	"codeberg.org/dergs/tidalwave/internal/g"
	"codeberg.org/dergs/tidalwave/internal/router"
	"codeberg.org/dergs/tidalwave/internal/ui/components/tracklist"
	"codeberg.org/dergs/tidalwave/pkg/schwifty/state"
	. "codeberg.org/dergs/tidalwave/pkg/schwifty/syntax"
	"codeberg.org/dergs/tidalwave/pkg/tidalapi"
	"codeberg.org/dergs/tidalwave/pkg/utils/imgutil"
	"github.com/infinytum/injector"
	"github.com/jwijenbergh/puregotk/v4/gdkpixbuf"
	"github.com/jwijenbergh/puregotk/v4/glib"
	"github.com/jwijenbergh/puregotk/v4/gtk"
)

func init() {
	router.Register("album", Album)
}

func Album(params router.Params) *router.Response {
	albumId, ok := params["id"].(string)
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

	coverState := state.NewStateful[*gdkpixbuf.Pixbuf](nil)
	for _, artwork := range album.Included.PlainArtworks(album.Data.Relationships.CoverArt.Data...) {
		go func() {
			pixbuf, err := injector.MustInject[*imgutil.ImgUtil]().LoadPixbuf(artwork.Attributes.Files.AtLeast(320).Href)
			if err != nil {
				slog.Error("failed to load album cover", "error", err)
				return
			}
			glib.IdleAddOnce(
				g.Ptr[glib.SourceOnceFunc](func(u uintptr) {
					coverState.SetValue(pixbuf)
					pixbuf.Unref()
				}),
				0,
			)
		}()
		break
	}

	list := tracklist.NewTrackList(
		"",
		tracklist.PositionColumn,
		tracklist.TitleColumn,
		tracklist.ArtistsColumn,
		tracklist.DurationColumn,
		tracklist.ButtonColumn,
		tracklist.ControlsColumn,
	)

	for _, track := range items.Included.Tracks(items.Data...) {
		list.AddTrack(&track)
	}

	return &router.Response{
		PageTitle: album.Data.Attributes.Title,
		View: VStack(
			HStack(
				AspectFrame(
					Image().
						BindPixbuf(coverState).
						PixelSize(146),
				).CornerRadius(10).Overflow(gtk.OverflowHiddenValue),
				VStack(
					Label(album.Data.Attributes.Title).
						FontSize(18).
						FontWeight(700).
						HAlign(gtk.AlignStartValue),
					Label(strings.Join(artists, ", ")).
						FontSize(16).
						FontWeight(500).
						HAlign(gtk.AlignStartValue),
					Label(album.Data.Attributes.ReleaseDate.Format("2006")).
						FontSize(16).
						FontWeight(500).
						HAlign(gtk.AlignStartValue),
					Label(fmt.Sprintf("%d Tracks (%s)", album.Data.Attributes.NumberOfItems, tidalapi.FormatDuration(int(album.Data.Attributes.Duration.Seconds())))).
						FontSize(14).
						FontWeight(600).
						HAlign(gtk.AlignStartValue).
						MarginTop(10),
				).MarginStart(20).VAlign(gtk.AlignCenterValue),
			),
			ScrolledWindow().
				Child(list).
				Policy(gtk.PolicyNeverValue, gtk.PolicyAutomaticValue).
				VExpand(true).
				MarginTop(20),
		).HMargin(40).VMargin(20),
	}
}
