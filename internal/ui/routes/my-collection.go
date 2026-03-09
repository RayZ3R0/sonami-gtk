package routes

import (
	"context"

	"github.com/RayZ3R0/sonami-gtk/internal/gettext"
	"github.com/RayZ3R0/sonami-gtk/internal/router"
	"github.com/RayZ3R0/sonami-gtk/internal/secrets"
	"github.com/RayZ3R0/sonami-gtk/internal/services/tidal/openapi"
	"github.com/RayZ3R0/sonami-gtk/internal/ui/components"
	"github.com/RayZ3R0/sonami-gtk/internal/ui/components/horizontal_list"
	"github.com/RayZ3R0/sonami-gtk/internal/ui/components/media_card"
	"github.com/RayZ3R0/sonami-gtk/internal/ui/components/tracklist"
	// "github.com/RayZ3R0/sonami-gtk/internal/ui/routes/my_collection"
	. "github.com/RayZ3R0/sonami-gtk/pkg/schwifty/syntax"
	"github.com/RayZ3R0/sonami-gtk/pkg/tidalapi"
	"codeberg.org/puregotk/puregotk/v4/gtk"
	"github.com/infinytum/injector"
)

func init() {
	// My Collection routes disabled in account-free mode — requires a real Tidal user account.
	// See hifi/deferred_features.md for details.
	// router.Register("my-collection", MyCollection)
	// router.Register("my-collection/albums", my_collection.Albums)
	// router.Register("my-collection/artists", my_collection.Artists)
	// router.Register("my-collection/playlists", my_collection.Playlists)
	// router.Register("my-collection/tracks", my_collection.Tracks)
}

func MyCollection() *router.Response {
	tidal := injector.MustInject[*tidalapi.TidalAPI]()
	userId := secrets.UserID()
	if userId == "" {
		return &router.Response{
			PageTitle: gettext.Get("My Collection"),
			View:      components.AuthRequired(gettext.Get("Please sign in to view your collection")),
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
		artistList.Append(media_card.NewArtist(openapi.NewArtistInfo(artist)))
	}

	albumList := horizontal_list.NewHorizontalList(gettext.Get("Albums")).SetPageMargin(40).SetViewAllRoute("my-collection/albums")
	for _, album := range userCollection.Included.Albums(userCollection.Data.Relationships.Albums.Data...) {
		albumList.Append(media_card.NewAlbum(openapi.NewAlbum(album)))
	}

	playlistList := horizontal_list.NewHorizontalList(gettext.Get("Playlists")).SetPageMargin(40).SetViewAllRoute("my-collection/playlists")
	for _, playlist := range userCollection.Included.Playlists(userCollection.Data.Relationships.Playlists.Data...) {
		playlistList.Append(media_card.NewPlaylist(openapi.NewPlaylist(playlist)))
	}

	trackList := tracklist.NewTrackList(
		tracklist.CoverColumn, tracklist.TitleAlbumColumn,
		tracklist.ArtistsColumn,
		tracklist.DurationColumn, tracklist.ControlsColumn,
	)
	for _, track := range userCollection.Included.Tracks(userCollection.Data.Relationships.Tracks.Data...) {
		trackList.AddTrack(openapi.NewTrack(track))
	}

	return &router.Response{
		PageTitle: gettext.Get("My Collection"),
		View: ScrolledWindow().
			Child(
				components.MainContent(
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
				),
			).
			Policy(gtk.PolicyNeverValue, gtk.PolicyAutomaticValue),
	}
}
