package routes

import (
	"context"

	"codeberg.org/dergs/tidalwave/internal/router"
	"codeberg.org/dergs/tidalwave/internal/secrets"
	"codeberg.org/dergs/tidalwave/internal/ui/components/horizontal_list"
	"codeberg.org/dergs/tidalwave/internal/ui/components/tracklist"
	. "codeberg.org/dergs/tidalwave/pkg/schwifty/syntax"
	"codeberg.org/dergs/tidalwave/pkg/tidalapi"
	"github.com/infinytum/injector"
	"github.com/jwijenbergh/puregotk/v4/gtk"
)

func init() {
	router.Register("my-collection", MyCollection)
}

func MyCollection(params router.Params) *router.Response {
	tidal := injector.MustInject[*tidalapi.TidalAPI]()
	userId := secrets.UserID()
	if userId == "" {
		return &router.Response{
			PageTitle: "My Collection",
			View:      Label("Please log in to view your collection"),
		}
	}

	userCollection, err := tidal.OpenAPI.V2.UserCollections.UserCollection(context.Background(), userId, "albums.coverArt", "artists.profileArt", "playlists.coverArt", "tracks.artists", "tracks.albums.coverArt")
	if err != nil {
		return &router.Response{
			PageTitle: "My Collection",
			View:      Label("Error loading collection"),
		}
	}

	artistList := horizontal_list.NewHorizontalList("Artists").SetPageMargin(40)
	for _, artist := range userCollection.Included.Artists(userCollection.Data.Relationships.Artists.Data...) {
		artistList.Append(horizontal_list.NewArtist(&artist))
	}

	albumList := horizontal_list.NewHorizontalList("Albums").SetPageMargin(40)
	for _, album := range userCollection.Included.Albums(userCollection.Data.Relationships.Albums.Data...) {
		albumList.Append(horizontal_list.NewAlbum(&album))
	}

	playlistList := horizontal_list.NewHorizontalList("Playlists").SetPageMargin(40)
	for _, playlist := range userCollection.Included.Playlists(userCollection.Data.Relationships.Playlists.Data...) {
		playlistList.Append(horizontal_list.NewPlaylist(&playlist))
	}

	trackList := tracklist.NewTrackList(
		"Tracks",
		tracklist.CoverColumn,
		tracklist.TitleAlbumColumn,
		tracklist.ArtistsColumn,
		tracklist.DurationColumn,
		tracklist.ButtonColumn,
		tracklist.ControlsColumn,
	)
	for _, track := range userCollection.Included.Tracks(userCollection.Data.Relationships.Tracks.Data...) {
		trackList.AddTrack(&track)
	}

	return &router.Response{
		PageTitle: "My Collection",
		View: ScrolledWindow().
			Child(
				VStack(
					artistList,
					albumList,
					playlistList,
					trackList.HMargin(40),
					Spacer(),
				).Spacing(25).VMargin(20).VAlign(gtk.AlignStartValue),
			).
			Policy(gtk.PolicyNeverValue, gtk.PolicyAutomaticValue),
	}
}
