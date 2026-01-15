package my_collection

import (
	"context"
	"log/slog"
	"unsafe"

	"codeberg.org/dergs/tonearm/internal/router"
	"codeberg.org/dergs/tonearm/internal/secrets"
	"codeberg.org/dergs/tonearm/internal/ui/components/tracklist"
	"codeberg.org/dergs/tonearm/pkg/schwifty"
	. "codeberg.org/dergs/tonearm/pkg/schwifty/syntax"
	"codeberg.org/dergs/tonearm/pkg/tidalapi"
	"codeberg.org/dergs/tonearm/pkg/tidalapi/models/openapi"
	"codeberg.org/dergs/tonearm/pkg/tidalapi/openapi/v2/user_collections"
	"codeberg.org/dergs/tonearm/pkg/tidalapi/pagination"
	"github.com/infinytum/injector"
	"github.com/jwijenbergh/puregotk/v4/gtk"
)

type ItemizedUserTracksCollection struct {
	*user_collections.UserCollections
}

func (i *ItemizedUserTracksCollection) Items(ctx context.Context, id, cursor string, include ...string) (*openapi.Response[[]openapi.Relationship], error) {
	return i.Tracks(ctx, id, cursor, include...)
}

func Tracks() *router.Response {
	tidal := injector.MustInject[*tidalapi.TidalAPI]()
	userId := secrets.UserID()
	if userId == "" {
		return &router.Response{
			PageTitle: "My Collection",
			View:      Label("Please log in to view your collection"),
		}
	}

	// userCollection, err := tidal.OpenAPI.V2.UserCollections.Tracks(context.Background(), userId, "", "tracks.artists", "tracks.albums.coverArt")
	paginatedCollection := &ItemizedUserTracksCollection{tidal.OpenAPI.V2.UserCollections}

	paginator := pagination.NewPaginator(paginatedCollection, userId, func(r *openapi.Response[[]openapi.Relationship]) []openapi.Track {
		return r.Included.Tracks(r.Data...)
	}, "tracks.artists", "tracks.albums.coverArt")

	userCollection, err := paginator.GetFirstPage()
	if err != nil {
		return &router.Response{
			PageTitle: "My Collection",
			Error:     err,
		}
	}

	trackList := tracklist.NewTrackList(
		tracklist.GroupedColumn(2, gtk.AlignStartValue, tracklist.CoverColumn, tracklist.TitleAlbumColumn),
		tracklist.ArtistsColumn,
		tracklist.ExpandButtonColumn(1),
		tracklist.GroupedColumn(1, gtk.AlignStartValue, tracklist.DurationColumn, tracklist.ControlsColumn),
	)

	for _, track := range userCollection {
		trackList.AddTrack(&track)
	}

	return &router.Response{
		PageTitle: "My Tracks",
		View: ScrolledWindow().
			Child(
				trackList.HMargin(40).VAlign(gtk.AlignStartValue),
			).
			ConnectEdgeReached(func(sw gtk.ScrolledWindow, pt gtk.PositionType) {
				if pt == gtk.PosBottomValue {
					go func() {
						if !paginator.IsConsumed() {
							items, err := paginator.Next()
							if err != nil {
								return
							}

							schwifty.OnMainThreadOnce(func(u uintptr) {
								var list *tracklist.TrackList[*openapi.Track]
								list = (*tracklist.TrackList[*openapi.Track])(unsafe.Pointer(u))
								for _, track := range items {
									list.AddTrack(&track)
								}
							}, uintptr(unsafe.Pointer(trackList)))
						} else {
							slog.Debug("No more tracks to fetch")
						}
					}()
				}
			}).
			Policy(gtk.PolicyNeverValue, gtk.PolicyAutomaticValue),
	}
}
