package shortcut_list

import (
	"strconv"

	"codeberg.org/dergs/tidalwave/pkg/schwifty"
	"codeberg.org/dergs/tidalwave/pkg/tidalapi"
	v2 "codeberg.org/dergs/tidalwave/pkg/tidalapi/models/v2"
	"github.com/jwijenbergh/puregotk/v4/glib"
)

func newArtist(id string, name string, coverUrl string) schwifty.Button {
	return NewShortcut(
		name,
		"",
		coverUrl,
	).ActionName("win.route.artist").ActionTargetValue(glib.NewVariantString(id))
}

func NewLegacyArtist(artist *v2.ArtistItemData) schwifty.Button {
	return newArtist(strconv.Itoa(artist.Id), artist.Name, tidalapi.ImageURL(artist.Picture))
}
