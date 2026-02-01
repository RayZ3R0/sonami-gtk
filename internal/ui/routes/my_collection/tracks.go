package my_collection

import (
	"codeberg.org/dergs/tonearm/internal/gettext"
	"codeberg.org/dergs/tonearm/internal/router"
	"codeberg.org/dergs/tonearm/internal/secrets"
	"codeberg.org/dergs/tonearm/internal/ui/components/tracklist"
	"codeberg.org/dergs/tonearm/internal/ui/pages"
	"codeberg.org/dergs/tonearm/pkg/schwifty"
	. "codeberg.org/dergs/tonearm/pkg/schwifty/syntax"
	"codeberg.org/dergs/tonearm/pkg/tidalapi"
	"codeberg.org/dergs/tonearm/pkg/tidalapi/models/openapi"
	"codeberg.org/dergs/tonearm/pkg/tidalapi/pagination"
	"github.com/infinytum/injector"
	"github.com/jwijenbergh/puregotk/v4/gtk"
)

func Tracks() *router.Response {
	tidal := injector.MustInject[*tidalapi.TidalAPI]()
	userId := secrets.UserID()
	if userId == "" {
		return &router.Response{
			PageTitle: gettext.Get("My Collection"),
			View:      Label(gettext.Get("Please log in to view your collection")),
		}
	}

	paginator := pagination.NewPaginator(tidal.OpenAPI.V2.UserCollections.Tracks, userId, func(r *openapi.Response[[]openapi.Relationship]) []openapi.Track {
		return r.Included.Tracks(r.Data...)
	}, "tracks.artists", "tracks.albums.coverArt")

	page, err := pages.NewPaginatedTracklistPage(paginator, func() *tracklist.TrackList[*openapi.Track] {
		return tracklist.NewTrackList(
			tracklist.GroupedColumn(2, gtk.AlignStartValue, tracklist.CoverColumn, tracklist.TitleAlbumColumn),
			tracklist.ArtistsColumn,
			tracklist.ExpandButtonColumn(1),
			tracklist.GroupedColumn(1, gtk.AlignEndValue, tracklist.DurationColumn, tracklist.ControlsColumn),
		)
	}, func(tl *tracklist.TrackList[*openapi.Track]) schwifty.BaseWidgetable {
		return tl.VPadding(20).HMargin(40).VAlign(gtk.AlignStartValue)
	})

	return &router.Response{
		PageTitle: gettext.Get("My Tracks"),
		Error:     err,
		View:      page,
	}
}
