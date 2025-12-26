package routes

import (
	"context"
	"errors"
	"fmt"

	"codeberg.org/dergs/tidalwave/internal/router"
	"codeberg.org/dergs/tidalwave/internal/ui/components/tracklist"
	. "codeberg.org/dergs/tidalwave/pkg/gui"
	"codeberg.org/dergs/tidalwave/pkg/tidalapi"
	"codeberg.org/dergs/tidalwave/pkg/tidalapi/models/openapi"
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

	// TODO: Handle pagination with scroll events
	items, err := tidal.OpenAPI.V2.Playlists.Items(context.Background(), playlistUUID, "", "items", "items.artists", "items.albums", "items.albums.coverArt")
	if err != nil {
		return router.FromError("Playlist", err)
	}

	creator := "TIDAL"
	for _, artist := range playlist.Included.PlainArtists(playlist.Data.Relationships.OwnerProfiles.Data...) {
		creator = artist.Attributes.Name
		break
	}

	cover := gtk.NewImage()
	cover.SetPixelSize(146)
	for _, artwork := range playlist.Included.PlainArtworks(playlist.Data.Relationships.CoverArt.Data...) {
		imgutil.AsyncGET(injector.MustInject[context.Context](), artwork.Attributes.Files.AtLeast(320).Href, imgutil.ImageSetterFromImage(cover))
		break
	}

	list := tracklist.NewTrackList(
		tracklist.CoverColumn,
		tracklist.TitleAlbumColumn,
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

	var playlistMetadata *TextImpl
	if playlist.Data.Attributes.PlaylistType != openapi.PlaylistTypeMix {
		playlistMetadata = Text(fmt.Sprintf("%d Tracks (%s)", playlist.Data.Attributes.NumberOfItems, tidalapi.FormatDuration(int(playlist.Data.Attributes.Duration.Seconds()))))
	} else {
		playlistMetadata = Text("Personal Mix")
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
					playlistMetadata.
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
