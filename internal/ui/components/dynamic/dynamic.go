package dynamic

import (
	"fmt"
	"strings"
	"time"

	"codeberg.org/dergs/tidalwave/internal/router"
	"codeberg.org/dergs/tidalwave/internal/ui/components"
	"codeberg.org/dergs/tidalwave/pkg/tidalapi"
	v2 "codeberg.org/dergs/tidalwave/pkg/tidalapi/models/v2"
	"github.com/diamondburned/gotk4/pkg/gtk/v4"
)

func ForItem(item v2.Item) gtk.Widgetter {
	switch item.Type {
	case v2.ItemTypeAlbum:
		card := components.NewMediaCard()
		card.SetTitle(item.Data.Album.Title)

		artists := make([]string, 0)
		for _, artist := range item.Data.Album.Artists {
			artists = append(artists, artist.Name)
		}

		releaseDate, _ := time.Parse(time.DateOnly, item.Data.Album.ReleaseDate)
		card.SetSubTitle(fmt.Sprintf("%s\n%s", strings.Join(artists, ", "), releaseDate.Format("2006")))
		card.LoadImage(tidalapi.ImageURL(item.Data.Album.Cover))
		card.ConnectClicked(func() {
			router.NavigateTo("album", router.Params{
				"id": item.Data.Album.Id,
			})
		})
		return card
	case v2.ItemTypePlaylist:
		card := components.NewMediaCard()
		card.SetTitle(item.Data.Playlist.Title)

		creator := "TIDAL"
		if item.Data.Playlist.Creator.Name != "" {
			creator = item.Data.Playlist.Creator.Name
		}

		card.SetSubTitle(fmt.Sprintf("%s\n%d Tracks", creator, item.Data.Playlist.NumberOfTracks))
		card.LoadImage(tidalapi.ImageURL(item.Data.Playlist.SquareImage))
		card.ConnectClicked(func() {
			router.NavigateTo("playlist", router.Params{
				"uuid": item.Data.Playlist.UUID,
			})
		})
		return card
	default:
		return gtk.NewLabel("Not implemented")
	}
}

func ForPageItem(item v2.PageItem) gtk.Widgetter {
	switch item.Type {
	case v2.ItemTypeHorizontalList:
		list := NewHorizontalList().SetTitle(item.Title).HMargin(40)
		for _, child := range item.Items {
			list.Append(ForItem(child))
		}
		return list
	case v2.ItemTypeTrackList:
		list := NewTrackList().SetTitle(item.Title)
		for i, track := range item.Items {
			data := track.Data.Track
			artists := make([]string, 0)
			for _, artist := range data.Artists {
				artists = append(artists, artist.Name)
			}

			parsedDuration := time.Second * time.Duration(data.Duration)

			// Track Lists are special and will always use a TrackListEntry
			trackListEntry := NewTrackListEntry(data.ID).
				SetAlbum(data.Album.Title).
				SetArtist(strings.Join(artists, ", ")).
				SetCoverFromURL(tidalapi.ImageURL(data.Album.Cover)).
				SetTime(parsedDuration.Round(time.Second).String()).
				SetTitle(data.Title)
			list.Append(trackListEntry, i)
		}
		return list.HMargin(40)
	default:
		return nil
	}
}
