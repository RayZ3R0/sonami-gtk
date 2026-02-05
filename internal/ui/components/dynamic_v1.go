package components

import (
	"strconv"
	"strings"
	"time"

	"codeberg.org/dergs/tonearm/internal/router"
	"codeberg.org/dergs/tonearm/internal/ui/components/horizontal_list"
	"codeberg.org/dergs/tonearm/internal/ui/components/media_card"
	"codeberg.org/dergs/tonearm/internal/ui/components/shortcut_list"
	"codeberg.org/dergs/tonearm/internal/ui/components/tracklist"
	"codeberg.org/dergs/tonearm/pkg/schwifty"
	. "codeberg.org/dergs/tonearm/pkg/schwifty/syntax"
	"codeberg.org/dergs/tonearm/pkg/tidalapi"
	v1 "codeberg.org/dergs/tonearm/pkg/tidalapi/models/v1"
	v2 "codeberg.org/dergs/tonearm/pkg/tidalapi/models/v2"
	"github.com/jwijenbergh/puregotk/v4/gtk"
)

func ForModule(module v1.Module) schwifty.BaseWidgetable {
	switch module.Type {
	case v1.ModuleTypeVideoList:
		// TODO: Implement video lists
		return HStack()
	case v1.ModuleTypeFeaturedPromotions:
		list := horizontal_list.NewHorizontalList(module.Title)
		for _, item := range module.Items {
			if item.Type == v1.ItemTypeCategoryPages {
				continue
			} else if item.Type == v1.ItemTypePlaylist {
				list.Append(media_card.NewMixGeneric(item.ArtifactID, item.ShortHeader, item.ShortSubHeader, tidalapi.ImageURLWithSize(item.ImageID, 550, 400)))
				continue
			} else {
				list.Append(HStack(
					Label("Unsupported\n"+string(item.Type)).
						Justify(gtk.JustifyCenterValue).
						Padding(30).HExpand(true).CornerRadius(10),
				).SizeRequest(192, -1).CSS("box:hover { background: alpha(var(--view-fg-color), 0.1); }").CornerRadius(10))
			}
		}
		return list.SetPageMargin(40)
	case v1.ModuleTypePlaylistList:
		list := horizontal_list.NewHorizontalList(module.Title)
		for _, item := range module.PagedList.Items {
			creator := "TIDAL"
			if len(item.Creators) > 0 {
				names := make([]string, len(item.Creators))
				for i, creator := range item.Creators {
					names[i] = creator.Name
				}
				creator = strings.Join(names, ", ")
			}
			list.Append(media_card.NewPlaylistGeneric(item.UUID, item.Title, creator, item.NumberOfTracks, tidalapi.ImageURL(item.SquareImage)))
			continue
		}
		return list.SetPageMargin(40)
	case v1.ModuleTypeArtistList:
		list := horizontal_list.NewHorizontalList(module.Title)
		for _, item := range module.PagedList.Items {
			list.Append(media_card.NewArtistGeneric(strconv.Itoa(item.ID), item.Name, tidalapi.ImageURL(item.Picture)))
			continue
		}
		return list.SetPageMargin(40)
	case v1.ModuleTypeAlbumList:
		list := horizontal_list.NewHorizontalList(module.Title)
		for _, item := range module.PagedList.Items {
			releaseDate, _ := time.Parse(time.DateOnly, item.ReleaseDate)
			artists := ""
			if len(item.Artists) > 0 {
				names := make([]string, len(item.Artists))
				for i, artist := range item.Artists {
					names[i] = artist.Name
				}
				artists = strings.Join(names, ", ")
			}
			list.Append(media_card.NewAlbumGeneric(strconv.Itoa(item.ID), item.Title, artists, releaseDate.Format("2006"), tidalapi.ImageURL(item.Cover)))
			continue
		}
		return list.SetPageMargin(40)
	case v1.ModuleTypeTrackList:
		list := tracklist.NewTrackList[*v2.TrackItemData](
			tracklist.GroupedColumn(2, gtk.AlignStartValue, tracklist.LegacyCoverColumn, tracklist.LegacyTitleAlbumColumn),
			tracklist.LegacyArtistsColumn,
			tracklist.LegacyExpandButtonColumn(1),
			tracklist.GroupedColumn(1, gtk.AlignEndValue, tracklist.LegacyDurationColumn, tracklist.LegacyControlsColumn),
		)
		for _, item := range module.PagedList.Items {
			artists := make([]v2.TrackItemDataArtist, 0)
			for _, artist := range item.Artists {
				artists = append(artists, v2.TrackItemDataArtist{
					ID:   artist.ID,
					Name: artist.Name,
				})
			}
			fakeTrack := &v2.TrackItemData{
				Album: v2.TrackItemDataAlbum{
					Cover: item.Album.Cover,
					ID:    item.Album.ID,
					Title: item.Album.Title,
				},
				Artists:   artists,
				Duration:  item.Duration,
				Following: false,
				ID:        item.ID,
				Title:     item.Title,
			}
			list.AddTrack(fakeTrack)
		}
		return VStack(
			NewRowTitle().SetTitle(module.Title),
			list,
		).HMargin(40)
	case v1.ModuleTypePageLinksCloud:
		list := shortcut_list.NewShortcutList()
		for _, item := range module.PagedList.Items {
			list.Append(shortcut_list.NewTextShortcut(item.Title, "").ConnectClicked(func(b gtk.Button) {
				router.Navigate(strings.ReplaceAll(item.APIPath, "pages/", "explore/"))
			}))
		}
		return VStack(
			NewRowTitle().SetTitle(module.Title).HPadding(40),
			list.HMargin(50),
		)
	case v1.ModuleTypePageLinks:
		list := shortcut_list.NewShortcutList()
		for _, item := range module.PagedList.Items {
			list.Append(shortcut_list.NewTextShortcut(item.Title, "").ConnectClicked(func(b gtk.Button) {
				router.Navigate(strings.ReplaceAll(item.APIPath, "pages/", "explore/"))
			}))
		}
		return VStack(
			NewRowTitle().SetTitle("More").HPadding(40),
			list.HMargin(50),
		)
	default:
		logger.Warn("Unsupported module type", "type", module.Type)
		return HStack(
			Label("Unsupported Element").
				Background("alpha(var(--view-fg-color), 0.1)").
				Padding(30).HExpand(true).CornerRadius(10),
		).HExpand(true).HMargin(40)
	}
}
