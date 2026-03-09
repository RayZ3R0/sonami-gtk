package search

import (
	"github.com/RayZ3R0/sonami-gtk/internal/gettext"
	"github.com/RayZ3R0/sonami-gtk/internal/services/tidal/openapi"
	"github.com/RayZ3R0/sonami-gtk/internal/ui/components"
	"github.com/RayZ3R0/sonami-gtk/internal/ui/components/horizontal_list"
	"github.com/RayZ3R0/sonami-gtk/internal/ui/components/media_card"
	"github.com/RayZ3R0/sonami-gtk/internal/ui/components/tracklist"
	"github.com/RayZ3R0/sonami-gtk/pkg/schwifty"
	. "github.com/RayZ3R0/sonami-gtk/pkg/schwifty/syntax"
	modelopenapi "github.com/RayZ3R0/sonami-gtk/pkg/tidalapi/models/openapi"
	"codeberg.org/puregotk/puregotk/v4/gtk"
)

func TopHits(searchResults *modelopenapi.SearchResult) schwifty.Box {
	artistList := horizontal_list.NewHorizontalList(gettext.Get("Artists")).SetPageMargin(40).SetViewAllRoute("search/" + searchResults.Data.ID + "/artists")
	for _, artist := range searchResults.Included.Artists(searchResults.Data.Relationships.TopHits.Data...) {
		artistList.Append(media_card.NewArtist(openapi.NewArtistInfo(artist)))
	}

	trackList := tracklist.NewTrackList(
		tracklist.CoverColumn, tracklist.TitleAlbumColumn,
		tracklist.ArtistsColumn,
		tracklist.DurationColumn, tracklist.ControlsColumn,
	)
	for _, track := range searchResults.Included.Tracks(searchResults.Data.Relationships.TopHits.Data...) {
		trackList.AddTrack(openapi.NewTrack(track))
	}

	albumList := horizontal_list.NewHorizontalList(gettext.Get("Albums")).SetPageMargin(40).SetViewAllRoute("search/" + searchResults.Data.ID + "/albums")
	for _, album := range searchResults.Included.Albums(searchResults.Data.Relationships.TopHits.Data...) {
		albumList.Append(media_card.NewAlbum(openapi.NewAlbum(album)))
	}

	playlistList := horizontal_list.NewHorizontalList(gettext.Get("Playlists")).SetPageMargin(40).SetViewAllRoute("search/" + searchResults.Data.ID + "/playlists")
	for _, playlist := range searchResults.Included.Playlists(searchResults.Data.Relationships.TopHits.Data...) {
		playlistList.Append(media_card.NewPlaylist(openapi.NewPlaylist(playlist)))
	}

	return VStack(
		artistList,
		albumList,
		playlistList,
		VStack(
			components.NewRowTitle().SetTitle(gettext.Get("Tracks")).SetViewAllRoute("search/"+searchResults.Data.ID+"/tracks"),
			trackList,
		).HMargin(40),
		Spacer(),
	).Spacing(25).VMargin(20).VAlign(gtk.AlignStartValue)
}
