package v2

type ItemType string

const (
	ItemTypeAlbum          ItemType = "ALBUM"
	ItemTypeDeepLink       ItemType = "DEEP_LINK"
	ItemTypeHorizontalList ItemType = "HORIZONTAL_LIST"
	ItemTypePlaylist       ItemType = "PLAYLIST"
	ItemTypeShortcutList   ItemType = "SHORTCUT_LIST"
	ItemTypeTrack          ItemType = "TRACK"
	ItemTypeTrackList      ItemType = "TRACK_LIST"
)
