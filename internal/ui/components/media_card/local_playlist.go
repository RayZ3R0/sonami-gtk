package media_card

import (
	"codeberg.org/puregotk/puregotk/v4/glib"
	"github.com/RayZ3R0/sonami-gtk/internal/gettext"
	"github.com/RayZ3R0/sonami-gtk/internal/localdb"
	"github.com/RayZ3R0/sonami-gtk/pkg/schwifty"
)

// NewLocalPlaylist builds a media card for a user-created local playlist.
func NewLocalPlaylist(p localdb.LocalPlaylist) schwifty.Button {
	return Card(
		p.Name,
		SubTitle(gettext.GetN("%d Track", "%d Tracks", p.TrackCount, p.TrackCount)),
		p.CoverURL,
	).ActionName("win.route.local-playlist").ActionTargetValue(glib.NewVariantString(p.ID))
}
