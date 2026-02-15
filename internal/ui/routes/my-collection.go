package routes

import (
	"context"

	"codeberg.org/dergs/tonearm/internal/gettext"
	"codeberg.org/dergs/tonearm/internal/router"
	"codeberg.org/dergs/tonearm/internal/secrets"
	"codeberg.org/dergs/tonearm/internal/ui/components"
	"codeberg.org/dergs/tonearm/internal/ui/components/horizontal_list"
	"codeberg.org/dergs/tonearm/internal/ui/components/media_card"
	"codeberg.org/dergs/tonearm/internal/ui/components/tracklist"
	"codeberg.org/dergs/tonearm/internal/ui/routes/my_collection"
	. "codeberg.org/dergs/tonearm/pkg/schwifty/syntax"
	"codeberg.org/dergs/tonearm/pkg/tidalapi"
	"github.com/infinytum/injector"
	"github.com/jwijenbergh/puregotk/v4/gtk"
)

func init() {
	router.Register("my-collection", MyCollection)
	router.Register("my-collection/albums", my_collection.Albums)
	router.Register("my-collection/artists", my_collection.Artists)
	router.Register("my-collection/playlists", my_collection.Playlists)
	router.Register("my-collection/tracks", my_collection.Tracks)
}

func MyCollection() *router.Response {
	tidal := injector.MustInject[*tidalapi.TidalAPI]()
	userId := secrets.UserID()
	if userId == "" {
		return &router.Response{
			PageTitle: gettext.Get("My Collection"),
			View: components.AuthRequired(gettext.Get("Please sign in to view your collection")),
		}
	}

	userCollection, err := tidal.OpenAPI.V2.UserCollections.UserCollection(context.Background(), userId, "albums.coverArt", "artists.profileArt", "playlists.coverArt", "tracks.artists", "tracks.albums.coverArt")
	if err != nil {
		return &router.Response{
			PageTitle: gettext.Get("My Collection"),
			View:      Label(gettext.Get("Error loading collection")),
		}
	}

	artistList := horizontal_list.NewHorizontalList(gettext.Get("Artists")).SetPageMargin(40).SetViewAllRoute("my-collection/artists")
	for _, artist := range userCollection.Included.Artists(userCollection.Data.Relationships.Artists.Data...) {
		artistList.Append(media_card.NewArtist(&artist))
	}

	albumList := horizontal_list.NewHorizontalList(gettext.Get("Albums")).SetPageMargin(40).SetViewAllRoute("my-collection/albums")
	for _, album := range userCollection.Included.Albums(userCollection.Data.Relationships.Albums.Data...) {
		albumList.Append(media_card.NewAlbum(&album))
	}

	playlistList := horizontal_list.NewHorizontalList(gettext.Get("Playlists")).SetPageMargin(40).SetViewAllRoute("my-collection/playlists")
	for _, playlist := range userCollection.Included.Playlists(userCollection.Data.Relationships.Playlists.Data...) {
		playlistList.Append(media_card.NewPlaylist(&playlist))
	}

	trackList := tracklist.NewTrackList(
		tracklist.GroupedColumn(2, gtk.AlignStartValue, tracklist.CoverColumn, tracklist.TitleAlbumColumn),
		tracklist.ArtistsColumn,
		tracklist.ExpandButtonColumn(1),
		tracklist.GroupedColumn(1, gtk.AlignEndValue, tracklist.DurationColumn, tracklist.ControlsColumn),
	)
	for _, track := range userCollection.Included.Tracks(userCollection.Data.Relationships.Tracks.Data...) {
		trackList.AddTrack(&track)
	}

	return &router.Response{
		PageTitle: gettext.Get("My Collection"),
		View: ScrolledWindow().
			Child(
				VStack(
					artistList,
					albumList,
					playlistList,
					VStack(
						components.NewRowTitle().SetTitle(gettext.Get("Tracks")).SetViewAllRoute("my-collection/tracks"),
						trackList,
					).HMargin(40),
					Spacer(),
				).Spacing(25).VMargin(20).VAlign(gtk.AlignStartValue),
			).
			Policy(gtk.PolicyNeverValue, gtk.PolicyAutomaticValue),
	}
}
