package resources

import (
	"codeberg.org/puregotk/puregotk/v4/gdk"
	"github.com/RayZ3R0/sonami-gtk/internal/g"
	"github.com/RayZ3R0/sonami-gtk/pkg/schwifty"
	"github.com/RayZ3R0/sonami-gtk/pkg/schwifty/tracking"
)

var MissingAlbum = g.Lazy(func() schwifty.Paintable {
	image := gdk.NewTextureFromResource("/io/github/rayz3r0/SonamiGtk/icons/hicolor/512x512/state/missing-album.png")
	image.Ref()
	tracking.SetFinalizer("Texture", image)
	return image
})
