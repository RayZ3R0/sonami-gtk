package my_collection

import (
	"context"

	"codeberg.org/dergs/tidalwave/internal/router"
	"codeberg.org/dergs/tidalwave/internal/secrets"
	"codeberg.org/dergs/tidalwave/internal/ui/components/tracklist"
	. "codeberg.org/dergs/tidalwave/pkg/schwifty/syntax"
	"codeberg.org/dergs/tidalwave/pkg/tidalapi"
	"github.com/infinytum/injector"
	"github.com/jwijenbergh/puregotk/v4/gtk"
)

func Tracks(params router.Params) *router.Response {
	tidal := injector.MustInject[*tidalapi.TidalAPI]()
	userId := secrets.UserID()
	if userId == "" {
		return &router.Response{
			PageTitle: "My Collection",
			View:      Label("Please log in to view your collection"),
		}
	}

	userCollection, err := tidal.OpenAPI.V2.UserCollections.Tracks(context.Background(), userId, "", "tracks.artists", "tracks.albums.coverArt")
	if err != nil {
		return &router.Response{
			PageTitle: "My Collection",
			Error:     err,
		}
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
	for _, track := range userCollection.Included.Tracks(userCollection.Data...) {
		trackList.AddTrack(&track)
	}

	return &router.Response{
		PageTitle: "My Collection",
		View: ScrolledWindow().
			Child(
				VStack(
					trackList.HMargin(40),
					Spacer(),
				).Spacing(25).VMargin(20).VAlign(gtk.AlignStartValue),
			).
			Policy(gtk.PolicyNeverValue, gtk.PolicyAutomaticValue),
	}
}
