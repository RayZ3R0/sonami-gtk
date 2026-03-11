package v2

type ItemType string

const (
	ItemTypeAlbum                     ItemType = "ALBUM"
	ItemTypeArtist                    ItemType = "ARTIST"
	ItemTypeArtistLink                ItemType = "ARTIST_LINK"
	ItemTypeDeepLink                  ItemType = "DEEP_LINK"
	ItemTypeHorizontalList            ItemType = "HORIZONTAL_LIST"
	ItemTypeHorizontalListWithContext ItemType = "HORIZONTAL_LIST_WITH_CONTEXT"
	ItemTypeLinksList                 ItemType = "LINKS_LIST"
	ItemTypeMix                       ItemType = "MIX"
	ItemTypePlaylist                  ItemType = "PLAYLIST"
	ItemTypeShortcutList              ItemType = "SHORTCUT_LIST"
	ItemTypeTrack                     ItemType = "TRACK"
	ItemTypeTrackCredits              ItemType = "TRACK_CREDITS"
	ItemTypeTrackList                 ItemType = "TRACK_LIST"
	ItemTypeVideo                     ItemType = "VIDEO"
	ItemTypeArtistTrackCreditsCard    ItemType = "ARTIST_TRACK_CREDITS_CARD"
)
