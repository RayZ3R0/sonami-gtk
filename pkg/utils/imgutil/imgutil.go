package imgutil

import (
	"codeberg.org/dergs/tonearm/pkg/schwifty"
	"codeberg.org/dergs/tonearm/pkg/utils/cacheutil"
	"github.com/jwijenbergh/puregotk/v4/gtk"
)

type ImgUtil struct {
	cache *cacheutil.Cache
}

func (i *ImgUtil) LoadIntoImage(url string, image *gtk.Image) {
	image.Ref()
	go func() {
		texture, err := i.Load(url)
		if err != nil {
			image.Unref()
			image = nil
			return
		}

		schwifty.OnMainThreadOnce(func(u uintptr) {
			image.SetFromPaintable(texture)
			texture.Unref()
			texture = nil
			image.Unref()
			image = nil
		}, 0)
	}()
}

func NewImgUtil(appId string) *ImgUtil {
	return &ImgUtil{
		cache: cacheutil.NewCache(appId, "images"),
	}
}
