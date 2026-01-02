package search

import (
	"codeberg.org/dergs/tidalwave/internal/ui/components/horizontal_list"
	"codeberg.org/dergs/tidalwave/internal/ui/components/media_card"
	"codeberg.org/dergs/tidalwave/internal/ui/components/tracklist"
	"codeberg.org/dergs/tidalwave/pkg/schwifty"
	. "codeberg.org/dergs/tidalwave/pkg/schwifty/syntax"
	"codeberg.org/dergs/tidalwave/pkg/tidalapi/models/openapi"
	"github.com/jwijenbergh/puregotk/v4/gtk"
)

func TopHits(searchResults *openapi.SearchResult) schwifty.Box {
	artistList := horizontal_list.NewHorizontalList("Artists").SetPageMargin(40)
	for _, artist := range searchResults.Included.Artists(searchResults.Data.Relationships.TopHits.Data...) {
		artistList.Append(media_card.NewArtist(&artist))
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
	for _, track := range searchResults.Included.Tracks(searchResults.Data.Relationships.TopHits.Data...) {
		trackList.AddTrack(&track)
	}

	albumList := horizontal_list.NewHorizontalList("Albums").SetPageMargin(40)
	for _, album := range searchResults.Included.Albums(searchResults.Data.Relationships.TopHits.Data...) {
		albumList.Append(media_card.NewAlbum(&album))
	}

	playlistList := horizontal_list.NewHorizontalList("Playlists").SetPageMargin(40)
	for _, playlist := range searchResults.Included.Playlists(searchResults.Data.Relationships.TopHits.Data...) {
		playlistList.Append(media_card.NewPlaylist(&playlist))
	}

	return VStack(
		artistList,
		trackList.HMargin(40),
		albumList,
		playlistList,
		Spacer(),
	).Spacing(25).VMargin(20).VAlign(gtk.AlignStartValue)
}
