package imgutil

import (
	"codeberg.org/dergs/tonearm/pkg/schwifty"
	"codeberg.org/dergs/tonearm/pkg/schwifty/tracking"
	"codeberg.org/dergs/tonearm/pkg/schwifty/utils/weak"
	"codeberg.org/dergs/tonearm/pkg/utils/cacheutil"
	"codeberg.org/puregotk/puregotk/v4/gtk"
)

type ImgUtil struct {
	cache *cacheutil.Cache
}

func (i *ImgUtil) LoadIntoImage(url string, image *gtk.Image) {
	ref := weak.NewWidgetRef(&image.Widget)
	go func() {
		texture, err := i.Load(url)
		if err != nil {
			return
		}

		schwifty.OnMainThreadOnce(func(u uintptr) {
			ref.Use(func(widget *gtk.Widget) {
				gtk.ImageNewFromInternalPtr(widget.GoPointer()).SetFromPaintable(texture)
			})
			texture = nil
		}, 0)
	}()
}

func (i *ImgUtil) LoadIntoImageCropped(url string, image *gtk.Image) {
	ref := weak.NewWidgetRef(&image.Widget)
	go func() {
		texture, err := i.Load(url)
		if err != nil {
			return
		}
		cropped := Crop(texture)
		texture = nil
		tracking.SetFinalizer("Texture", cropped)

		schwifty.OnMainThreadOnce(func(u uintptr) {
			ref.Use(func(widget *gtk.Widget) {
				gtk.ImageNewFromInternalPtr(widget.GoPointer()).SetFromPaintable(cropped)
			})
			cropped = nil
		}, 0)
	}()
}

func NewImgUtil(appId string) *ImgUtil {
	return &ImgUtil{
		cache: cacheutil.NewCache(appId, "images"),
	}
}
