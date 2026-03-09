package media_card

import (
	"strings"

	"github.com/RayZ3R0/sonami-gtk/pkg/schwifty"
	. "github.com/RayZ3R0/sonami-gtk/pkg/schwifty/syntax"
	"github.com/RayZ3R0/sonami-gtk/pkg/sonami"
	"codeberg.org/puregotk/puregotk/v4/glib"
)

func NewAlbum(album sonami.Album) schwifty.Button {
	return Card(
		album.Title(),
		VStack(
			SubTitle(strings.Join(album.Artists().Names(), ", ")),
			SubTitle(album.ReleasedAt().Format("2006")),
		),
		album.Cover(172),
	).ActionName("win.route.album").ActionTargetValue(glib.NewVariantString(album.ID()))
}
