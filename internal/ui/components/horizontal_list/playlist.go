package horizontal_list

import (
	"fmt"

	"codeberg.org/dergs/tidalwave/pkg/schwifty"
	. "codeberg.org/dergs/tidalwave/pkg/schwifty/syntax"
	"codeberg.org/dergs/tidalwave/pkg/tidalapi"
	v2 "codeberg.org/dergs/tidalwave/pkg/tidalapi/models/v2"
	"github.com/jwijenbergh/puregotk/v4/glib"
)

func newPlaylist(id string, title string, creator string, itemCount int, coverUrl string) schwifty.Button {
	return Card(
		title,
		VStack(
			SubTitle(creator),
			SubTitle(fmt.Sprintf("%d Tracks", itemCount)),
		),
		coverUrl,
	).ActionName("win.route.playlist").ActionTargetValue(glib.NewVariantString(id))
}

func NewPlaylist(playlist *v2.PlaylistItemData) schwifty.Button {
	creator := "TIDAL"
	if playlist.Creator.Name != "" {
		creator = playlist.Creator.Name
	}
	return newPlaylist(playlist.UUID, playlist.Title, creator, playlist.NumberOfTracks, tidalapi.ImageURL(playlist.SquareImage))
}
