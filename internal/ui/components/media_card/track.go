package media_card

import (
	"strings"

	"codeberg.org/puregotk/puregotk/v4/glib"
	"github.com/RayZ3R0/sonami-gtk/pkg/schwifty"
	"github.com/RayZ3R0/sonami-gtk/pkg/sonami"
)

func NewTrack(track sonami.Track) schwifty.Button {
	return Card(
		track.Title(),
		SubTitle(strings.Join(track.Artists().Names(), ", ")),
		track.Cover(192),
	).ActionName("win.player.play-track").ActionTargetValue(glib.NewVariantString(track.ID()))
}
