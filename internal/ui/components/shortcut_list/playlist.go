package shortcut_list

import (
	"fmt"

	"codeberg.org/dergs/tidalwave/pkg/schwifty"
	"codeberg.org/dergs/tidalwave/pkg/tidalapi"
	v2 "codeberg.org/dergs/tidalwave/pkg/tidalapi/models/v2"
	"github.com/jwijenbergh/puregotk/v4/glib"
)

func newPlaylist(id string, title string, itemCount int, coverUrl string) schwifty.Button {
	return NewShortcut(
		title,
		fmt.Sprintf("%d Tracks", itemCount),
		coverUrl,
	).ActionName("win.route.playlist").ActionTargetValue(glib.NewVariantString(id))
}

func NewLegacyPlaylist(playlist *v2.PlaylistItemData) schwifty.Button {
	return newPlaylist(playlist.UUID, playlist.Title, playlist.NumberOfTracks, tidalapi.ImageURL(playlist.SquareImage))
}
