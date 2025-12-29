package components

import (
	"fmt"
	"strings"

	"codeberg.org/dergs/tidalwave/internal/ui/components/horizontal_list"
	"codeberg.org/dergs/tidalwave/internal/ui/components/shortcut_list"
	"codeberg.org/dergs/tidalwave/internal/ui/components/tracklist"
	"codeberg.org/dergs/tidalwave/pkg/schwifty"
	"codeberg.org/dergs/tidalwave/pkg/schwifty/syntax"
	. "codeberg.org/dergs/tidalwave/pkg/schwifty/syntax"
	"codeberg.org/dergs/tidalwave/pkg/tidalapi"
	v2 "codeberg.org/dergs/tidalwave/pkg/tidalapi/models/v2"
)

func ForPageItem(pageItem v2.PageItem) schwifty.BaseWidgetable {
	switch pageItem.Type {
	case v2.ItemTypeHorizontalList:
		list := horizontal_list.NewHorizontalList(pageItem.Title)
		for _, item := range pageItem.Items {
			if item.Type == v2.ItemTypeAlbum {
				album := item.Data.Album
				list.Append(horizontal_list.NewLegacyAlbum(album))
			} else if item.Type == v2.ItemTypePlaylist {
				playlist := item.Data.Playlist
				list.Append(horizontal_list.NewPlaylist(playlist))
			} else if item.Type == v2.ItemTypeArtist {
				artist := item.Data.Artist
				list.Append(horizontal_list.NewLegacyArtist(artist))
			} else if item.Type == v2.ItemTypeMix {
				mix := item.Data.Mix
				list.Append(horizontal_list.NewLegacyMix(mix))
			} else {
				list.Append(syntax.Label("Unsupported: " + string(item.Type)).HMargin(10))
			}
		}
		return list.SetPageMargin(40)
	case v2.ItemTypeTrackList:
		list := tracklist.NewLegacyTrackList(
			pageItem.Title,
			tracklist.LegacyCoverColumn,
			tracklist.LegacyTitleAlbumColumn,
			tracklist.LegacyArtistsColumn,
			tracklist.LegacyDurationColumn,
			tracklist.LegacyButtonColumn,
			tracklist.LegacyControlsColumn,
		)
		for _, track := range pageItem.Items {
			list.AddLegacyTrack(track.Data.Track)
		}
		return list.HMargin(40)
	case v2.ItemTypeShortcutList:
		list := shortcut_list.NewShortcutList()
		for _, item := range pageItem.Items {
			if item.Type == v2.ItemTypeDeepLink {
				deeplink := item.Data.DeepLink
				list.Append(shortcut_list.NewShortcut(deeplink.Title, "", ""))
			} else if item.Type == v2.ItemTypeAlbum {
				album := item.Data.Album
				artists := make([]string, 0)
				for _, artist := range album.Artists {
					artists = append(artists, artist.Name)
				}
				list.Append(shortcut_list.NewShortcut(album.Title, strings.Join(artists, ", "), tidalapi.ImageURL(album.Cover)))
			} else if item.Type == v2.ItemTypeArtist {
				artist := item.Data.Artist
				list.Append(shortcut_list.NewShortcut(artist.Name, "", tidalapi.ImageURL(artist.Picture)))
			} else if item.Type == v2.ItemTypePlaylist {
				playlist := item.Data.Playlist
				list.Append(shortcut_list.NewShortcut(playlist.Title, fmt.Sprintf("%d Tracks", playlist.NumberOfTracks), tidalapi.ImageURL(playlist.SquareImage)))
			} else if item.Type == v2.ItemTypeMix {
				mix := item.Data.Mix
				list.Append(shortcut_list.NewShortcut(mix.TitleTextInfo.Text, mix.SubtitleTextInfo.Text, mix.MixImages[0].URL))
			} else {
				list.Append(syntax.Label("Unsupported: " + string(item.Type)).HMargin(10))
			}
		}
		return list.HMargin(40)
	default:
		return HStack(
			Label("Unsupported Element").
				Background("alpha(var(--view-fg-color), 0.1)").
				Padding(30).HExpand(true).CornerRadius(10),
		).HExpand(true).HMargin(40)
	}
}
