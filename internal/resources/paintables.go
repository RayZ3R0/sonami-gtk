package resources

import (
	"codeberg.org/dergs/tonearm/pkg/schwifty"
	"github.com/jwijenbergh/puregotk/v4/gdk"
)

var MissingAlbum = func() schwifty.Paintable {
	image := gdk.NewTextureFromResource("/dev/dergs/Tonearm/icons/hicolor/512x512/state/missing-album.png")
	return image
}
