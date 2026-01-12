package resources

import (
	"codeberg.org/dergs/tidalwave/internal/g"
	"codeberg.org/dergs/tidalwave/pkg/schwifty"
	"github.com/jwijenbergh/puregotk/v4/gdk"
)

var missingAlbum = g.Lazy(func() schwifty.Paintable {
	return gdk.NewTextureFromResource("/dev/dergs/tidalwave/icons/scalable/state/missing-album.svg")
})

var MissingAlbum = g.Lazy(func() schwifty.Paintable {
	image := missingAlbum()
	image.Ref()
	return image
})
