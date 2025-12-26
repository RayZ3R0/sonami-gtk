package search

import (
	"fmt"
	"strconv"
	"strings"

	"codeberg.org/dergs/tidalwave/internal/router"
	"codeberg.org/dergs/tidalwave/internal/ui/components"
	"codeberg.org/dergs/tidalwave/internal/ui/components/dynamic"
	"codeberg.org/dergs/tidalwave/internal/ui/components/tracklist"
	"codeberg.org/dergs/tidalwave/pkg/gui"
	"codeberg.org/dergs/tidalwave/pkg/tidalapi/models/openapi"
	"github.com/diamondburned/gotk4/pkg/gtk/v4"
)

func TopHits(res *openapi.SearchResult) gtk.Widgetter {
	artistList := dynamic.NewHorizontalList().SetTitle("Artists").HMargin(40)
	artists := res.Included.Artists(res.Data.Relationships.TopHits.Data...)
	for i := 0; i < len(artists); i += 2 {
		stack := gui.VStack(forArtist(&artists[i])).Spacing(10).HMargin(5)

		if i+1 < len(artists) {
			stack.Append(forArtist(&artists[i+1]))
		}

		artistList.Append(stack)
	}

	trackList := tracklist.NewTrackList(
		tracklist.CoverColumn,
		tracklist.TitleAlbumColumn,
		tracklist.ArtistsColumn,
		tracklist.DurationColumn,
		tracklist.BoxColumn,
		tracklist.ControlsColumn,
	).SetTitle("Tracks")

	for _, track := range res.Included.Tracks(res.Data.Relationships.TopHits.Data...) {
		trackList.AddTrack(&track)
	}

	albumList := dynamic.NewHorizontalList().SetTitle("Albums").HMargin(40)
	for _, album := range res.Included.Albums(res.Data.Relationships.TopHits.Data...) {
		card := components.NewMediaCard()
		card.SetTitle(album.Data.Attributes.Title)

		artists := make([]string, 0)
		for _, artist := range album.Included.PlainArtists(album.Data.Relationships.Artists.Data...) {
			artists = append(artists, artist.Attributes.Name)
		}

		for _, artwork := range album.Included.PlainArtworks(album.Data.Relationships.CoverArt.Data...) {
			card.LoadImage(artwork.Attributes.Files.AtLeast(320).Href)
			break
		}

		card.SetSubTitle(fmt.Sprintf("%s\n%s", strings.Join(artists, ", "), album.Data.Attributes.ReleaseDate.Format("2006")))
		card.ConnectClicked(func() {
			id, _ := strconv.Atoi(album.Data.ID)
			router.Navigate("album", router.Params{
				"id": id,
			})
		})
		albumList.Append(card)
	}

	playlistList := dynamic.NewHorizontalList().SetTitle("Playlists").HMargin(40)
	for _, playlist := range res.Included.Playlists(res.Data.Relationships.TopHits.Data...) {
		card := components.NewMediaCard()
		card.SetTitle(playlist.Data.Attributes.Name)

		creator := "TIDAL"
		for _, profile := range playlist.Included.PlainArtists(playlist.Data.Relationships.OwnerProfiles.Data...) {
			creator = profile.Attributes.Name
			break
		}

		for _, artwork := range playlist.Included.PlainArtworks(playlist.Data.Relationships.CoverArt.Data...) {
			card.LoadImage(artwork.Attributes.Files.AtLeast(320).Href)
			break
		}

		card.SetSubTitle(fmt.Sprintf("%s\n%d Tracks", creator, playlist.Data.Attributes.NumberOfItems))
		card.ConnectClicked(func() {
			router.Navigate("playlist", router.Params{
				"uuid": playlist.Data.ID,
			})
		})
		playlistList.Append(card)
	}

	return gui.VStack(
		artistList,
		trackList.HMargin(40),
		albumList,
		playlistList,
		gui.Spacer(),
	).Spacing(25).VMargin(20)
}

func forArtist(artist *openapi.Artist) gtk.Widgetter {
	shortcut := components.NewShortcut()
	shortcut.SetTitle(artist.Data.Attributes.Name)
	for _, artwork := range artist.Included.PlainArtworks(artist.Data.Relationship.ProfileArt.Data...) {
		shortcut.LoadCover(artwork.Attributes.Files.AtLeast(320).Href)
		break
	}
	return shortcut
}
