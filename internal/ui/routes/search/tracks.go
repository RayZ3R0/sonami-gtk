package search

import (
	"codeberg.org/dergs/tonearm/internal/gettext"
	"codeberg.org/dergs/tonearm/internal/router"
	"codeberg.org/dergs/tonearm/internal/services/tidal/openapi"
	"codeberg.org/dergs/tonearm/internal/ui/components/tracklist"
	"codeberg.org/dergs/tonearm/internal/ui/pages"
	"codeberg.org/dergs/tonearm/pkg/schwifty"
	"codeberg.org/dergs/tonearm/pkg/tidalapi"
	modelopenapi "codeberg.org/dergs/tonearm/pkg/tidalapi/models/openapi"
	"codeberg.org/dergs/tonearm/pkg/tonearm"
	"github.com/infinytum/injector"
	"github.com/jwijenbergh/puregotk/v4/gtk"
)

func init() {
	router.Register("search/:query/tracks", tracks)
}

func tracks(query string) *router.Response {
	tidal := injector.MustInject[*tidalapi.TidalAPI]()

	paginator := openapi.NewPaginator(tidal.OpenAPI.V2.SearchResults.Tracks, query, func(r *modelopenapi.Response[[]modelopenapi.Relationship]) []tonearm.Track {
		results := r.Included.Tracks(r.Data...)
		tracks := make([]tonearm.Track, len(results))
		for i, track := range results {
			tracks[i] = openapi.NewTrack(track)
		}
		return tracks
	}, "tracks.albums.coverArt", "tracks.albums.artists")

	page, err := pages.NewPaginatedTracklistPage(paginator, func() *tracklist.TrackList {
		return tracklist.NewTrackList(
			tracklist.GroupedColumn(2, gtk.AlignStartValue, tracklist.CoverColumn, tracklist.TitleAlbumColumn),
			tracklist.ArtistsColumn,
			tracklist.ExpandButtonColumn(1),
			tracklist.GroupedColumn(1, gtk.AlignEndValue, tracklist.DurationColumn, tracklist.ControlsColumn),
		)
	}, func(tl *tracklist.TrackList) schwifty.BaseWidgetable {
		return tl.VPadding(20).HMargin(40).VAlign(gtk.AlignStartValue)
	})

	return &router.Response{
		PageTitle: gettext.Get("Search"),
		Error:     err,
		View:      page,
	}
}
