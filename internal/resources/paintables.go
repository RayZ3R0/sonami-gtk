package resources

import (
	"codeberg.org/dergs/tidalwave/internal/g"
	"github.com/jwijenbergh/puregotk/v4/gdk"
)

var MissingAlbum = g.Lazy(func() gdk.Paintable {
	return gdk.NewTextureFromResource("/org/codeberg/dergs/tidalwave/icons/scalable/state/missing-album.svg")
})
