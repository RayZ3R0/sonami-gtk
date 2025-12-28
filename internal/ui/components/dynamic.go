package components

import (
	"codeberg.org/dergs/tidalwave/internal/ui/components/tracklist"
	"codeberg.org/dergs/tidalwave/pkg/schwifty"
	. "codeberg.org/dergs/tidalwave/pkg/schwifty/syntax"
	v2 "codeberg.org/dergs/tidalwave/pkg/tidalapi/models/v2"
)

func ForPageItem(pageItem v2.PageItem) schwifty.BaseWidgetable {
	switch pageItem.Type {
	case v2.ItemTypeTrackList:
		list := tracklist.NewLegacyTrackList(
			tracklist.LegacyCoverColumn,
			tracklist.LegacyTitleAlbumColumn,
			tracklist.LegacyArtistsColumn,
			tracklist.LegacyDurationColumn,
			tracklist.LegacyButtonColumn,
			tracklist.LegacyControlsColumn,
		).SetTitle(pageItem.Title)
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
