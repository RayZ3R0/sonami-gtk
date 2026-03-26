package shortcut_list

import (
	"strconv"

	"codeberg.org/puregotk/puregotk/v4/glib"
	"github.com/RayZ3R0/sonami-gtk/pkg/schwifty"
	"github.com/RayZ3R0/sonami-gtk/pkg/tidalapi"
	v2 "github.com/RayZ3R0/sonami-gtk/pkg/tidalapi/models/v2"
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
