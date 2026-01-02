package components

import (
	"codeberg.org/dergs/tidalwave/internal/ui/components/horizontal_list"
	"codeberg.org/dergs/tidalwave/internal/ui/components/media_card"
	"codeberg.org/dergs/tidalwave/internal/ui/components/shortcut_list"
	"codeberg.org/dergs/tidalwave/internal/ui/components/tracklist"
	"codeberg.org/dergs/tidalwave/pkg/schwifty"
	"codeberg.org/dergs/tidalwave/pkg/schwifty/syntax"
	. "codeberg.org/dergs/tidalwave/pkg/schwifty/syntax"
	v2 "codeberg.org/dergs/tidalwave/pkg/tidalapi/models/v2"
	"github.com/jwijenbergh/puregotk/v4/gtk"
)

func ForPageItem(pageItem v2.PageItem) schwifty.BaseWidgetable {
	switch pageItem.Type {
	case v2.ItemTypeHorizontalList:
		list := horizontal_list.NewHorizontalList(pageItem.Title)
		for _, item := range pageItem.Items {
			if item.Type == v2.ItemTypeAlbum {
				list.Append(media_card.NewLegacyAlbum(item.Data.Album))
			} else if item.Type == v2.ItemTypePlaylist {
				list.Append(media_card.NewLegacyPlaylist(item.Data.Playlist))
			} else if item.Type == v2.ItemTypeArtist {
				list.Append(media_card.NewLegacyArtist(item.Data.Artist))
			} else if item.Type == v2.ItemTypeMix {
				list.Append(media_card.NewLegacyMix(item.Data.Mix))
			} else if item.Type == v2.ItemTypeTrack {
				list.Append(media_card.NewLegacyTrack(item.Data.Track))
			} else if item.Type == v2.ItemTypeDeepLink {
				list.Append(media_card.NewLegacyDeeplink(item.Data.DeepLink))
			} else {
				list.Append(HStack(
					Label("Unsupported\n"+string(item.Type)).
						Justify(gtk.JustifyCenterValue).
						Padding(30).HExpand(true).CornerRadius(10),
				).SizeRequest(192, -1).CSS("box:hover { background: alpha(var(--view-fg-color), 0.1) }").CornerRadius(10))
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
				list.Append(shortcut_list.NewLegacyDeepLink(item.Data.DeepLink))
			} else if item.Type == v2.ItemTypeAlbum {
				list.Append(shortcut_list.NewLegacyAlbum(item.Data.Album))
			} else if item.Type == v2.ItemTypeArtist {
				list.Append(shortcut_list.NewLegacyArtist(item.Data.Artist))
			} else if item.Type == v2.ItemTypePlaylist {
				list.Append(shortcut_list.NewLegacyPlaylist(item.Data.Playlist))
			} else if item.Type == v2.ItemTypeMix {
				list.Append(shortcut_list.NewLegacyMix(item.Data.Mix))
			} else {
				list.Append(syntax.Label("Unsupported: " + string(item.Type)).HMargin(10))
			}
		}
		return list.HMargin(50)
	default:
		return HStack(
			Label("Unsupported Element").
				Background("alpha(var(--view-fg-color), 0.1)").
				Padding(30).HExpand(true).CornerRadius(10),
		).HExpand(true).HMargin(40)
	}
}
