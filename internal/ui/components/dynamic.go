package components

import (
	"codeberg.org/dergs/tidalwave/internal/ui/components/horizontal_list"
	"codeberg.org/dergs/tidalwave/internal/ui/components/tracklist"
	"codeberg.org/dergs/tidalwave/pkg/schwifty"
	"codeberg.org/dergs/tidalwave/pkg/schwifty/syntax"
	. "codeberg.org/dergs/tidalwave/pkg/schwifty/syntax"
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
			} else {
				list.Append(syntax.Label(string(item.Type)).HMargin(10))
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
	default:
		return HStack(
			Spacer().VExpand(false),
			Label("Unsupported Element").VMargin(30),
			Spacer().VExpand(false),
		)
	}
}
