package dynamic

import (
	"fmt"
	"strings"
	"time"

	"codeberg.org/dergs/tidalwave/internal/router"
	"codeberg.org/dergs/tidalwave/internal/ui/components"
	"codeberg.org/dergs/tidalwave/internal/ui/components/tracklist"
	"codeberg.org/dergs/tidalwave/pkg/gui"
	"codeberg.org/dergs/tidalwave/pkg/tidalapi"
	v2 "codeberg.org/dergs/tidalwave/pkg/tidalapi/models/v2"
	"github.com/diamondburned/gotk4/pkg/gtk/v4"
)

func ForLegacyItem(item v2.Item) gtk.Widgetter {
	switch item.Type {
	case v2.ItemTypeAlbum:
		card := components.NewMediaCard()
		card.SetTitle(item.Data.Album.Title)

		artists := make([]string, 0)
		for _, artist := range item.Data.Album.Artists {
			artists = append(artists, artist.Name)
		}

		releaseDate, _ := time.Parse(time.DateOnly, item.Data.Album.ReleaseDate)
		card.SetSubTitle(fmt.Sprintf("%s\n%s", strings.Join(artists, ", "), releaseDate.Format("2006")), 1)
		card.LoadImage(tidalapi.ImageURL(item.Data.Album.Cover))
		card.ConnectClicked(func() {
			router.Navigate("album", router.Params{
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

		card.SetSubTitle(fmt.Sprintf("%s\n%d Tracks", creator, item.Data.Playlist.NumberOfTracks), 1)
		card.LoadImage(tidalapi.ImageURL(item.Data.Playlist.SquareImage))
		card.ConnectClicked(func() {
			router.Navigate("playlist", router.Params{
				"uuid": item.Data.Playlist.UUID,
			})
		})
		return card
	case v2.ItemTypeMix:
		card := components.NewMediaCard()
		card.SetTitle(item.Data.Mix.TitleTextInfo.Text)
		card.SetSubTitle(item.Data.Mix.SubtitleTextInfo.Text, 2)
		card.LoadImage(item.Data.Mix.MixImages[0].URL)
		card.ConnectClicked(func() {
			router.Navigate("playlist", router.Params{
				"uuid": item.Data.Mix.Id,
			})
		})
		return card
	case v2.ItemTypeArtist:
		card := components.NewMediaCard()
		card.SetTitle(item.Data.Artist.Name)
		card.LoadImage(tidalapi.ImageURL(item.Data.Artist.Picture))
		card.ConnectClicked(func() {
			router.Navigate("artist", router.Params{
				"id": item.Data.Artist.Id,
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
			list.Append(ForLegacyItem(child))
		}
		return list
	case v2.ItemTypeTrackList:
		list := tracklist.NewLegacyTrackList(
			tracklist.LegacyCoverColumn,
			tracklist.LegacyTitleAlbumColumn,
			tracklist.LegacyArtistsColumn,
			tracklist.LegacyDurationColumn,
			tracklist.LegacyBoxColumn,
			tracklist.LegacyControlsColumn,
		).SetTitle(item.Title)
		for _, track := range item.Items {
			list.AddLegacyTrack(track.Data.Track)
		}
		return list.HMargin(40)
	case v2.ItemTypeShortcutList:
		list := NewHorizontalList().SetTitle(item.Title).HMargin(40)
		for i := 0; i < len(item.Items); i += 2 {
			stack := gui.VStack(forShortcut(item.Items[i])).Spacing(10).HMargin(5)

			if i+1 < len(item.Items) {
				second := item.Items[i+1]
				stack.Append(forShortcut(second))
			}

			list.Append(stack)
		}
		return list
	default:
		return nil
	}
}

func forShortcut(child v2.Item) gtk.Widgetter {
	shortcut := components.NewShortcut()
	switch child.Type {
	case v2.ItemTypeDeepLink:
		shortcut.SetTitle(child.Data.DeepLink.Title)
	default:
		shortcut.SetTitle("Unsupported")
		shortcut.SetSubTitle("Element not supported")
	}
	return shortcut
}
