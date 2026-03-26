package media_card

import (
	"codeberg.org/puregotk/puregotk/v4/glib"
	"github.com/RayZ3R0/sonami-gtk/internal/gettext"
	"github.com/RayZ3R0/sonami-gtk/pkg/schwifty"
	. "github.com/RayZ3R0/sonami-gtk/pkg/schwifty/syntax"
	"github.com/RayZ3R0/sonami-gtk/pkg/sonami"
)

func NewPlaylist(playlist sonami.Playlist) schwifty.Button {
	creatorName := "TIDAL"
	if creator := playlist.Creator(); creator != nil {
		creatorName = creator.Title()
	}
	return Card(
		playlist.Title(),
		VStack(
			SubTitle(creatorName),
			SubTitle(gettext.GetN("%d Track", "%d Tracks", playlist.Count(), playlist.Count())),
		),
		playlist.Cover(192),
	).ActionName("win.route.playlist").ActionTargetValue(glib.NewVariantString(playlist.ID()))
}
