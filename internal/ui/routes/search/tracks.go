package search

import (
	"codeberg.org/puregotk/puregotk/v4/gtk"
	"github.com/RayZ3R0/sonami-gtk/internal/gettext"
	"github.com/RayZ3R0/sonami-gtk/internal/router"
	"github.com/RayZ3R0/sonami-gtk/internal/services/tidal/openapi"
	"github.com/RayZ3R0/sonami-gtk/internal/ui/components/tracklist"
	"github.com/RayZ3R0/sonami-gtk/internal/ui/pages"
	"github.com/RayZ3R0/sonami-gtk/pkg/schwifty"
	"github.com/RayZ3R0/sonami-gtk/pkg/sonami"
	"github.com/RayZ3R0/sonami-gtk/pkg/tidalapi"
	modelopenapi "github.com/RayZ3R0/sonami-gtk/pkg/tidalapi/models/openapi"
	"github.com/infinytum/injector"
)

func init() {
	router.Register("search/:query/tracks", tracks)
}

func tracks(query string) *router.Response {
	tidal := injector.MustInject[*tidalapi.TidalAPI]()

	paginator := openapi.NewPaginator(tidal.OpenAPI.V2.SearchResults.Tracks, query, func(r *modelopenapi.Response[[]modelopenapi.Relationship]) []sonami.Track {
		results := r.Included.Tracks(r.Data...)
		tracks := make([]sonami.Track, len(results))
		for i, track := range results {
			tracks[i] = openapi.NewTrack(track)
		}
		return tracks
	}, "tracks.albums.coverArt", "tracks.albums.artists")

	page, err := pages.NewPaginatedTracklistPage(paginator, func(tl *tracklist.TrackList) schwifty.BaseWidgetable {
		return tl.HMargin(40).VAlign(gtk.AlignStartValue)
	}, tracklist.CoverColumn, tracklist.TitleAlbumColumn, tracklist.ArtistsColumn, tracklist.DurationColumn, tracklist.ControlsColumn)

	return &router.Response{
		PageTitle: gettext.Get("Search"),
		Error:     err,
		View:      page,
	}
}
