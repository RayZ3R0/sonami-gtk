package media_card

import (
	"github.com/RayZ3R0/sonami-gtk/pkg/schwifty"
	. "github.com/RayZ3R0/sonami-gtk/pkg/schwifty/syntax"
	"github.com/RayZ3R0/sonami-gtk/pkg/sonami"
	"codeberg.org/puregotk/puregotk/v4/glib"
)

func NewArtist(artist sonami.ArtistInfo) schwifty.Button {
	return Card(
		artist.Title(),
		HStack(),
		artist.Cover(192),
	).ActionName("win.route.artist").ActionTargetValue(glib.NewVariantString(artist.ID()))
}
