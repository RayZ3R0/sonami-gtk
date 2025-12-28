package imgutil

import (
	"codeberg.org/dergs/tidalwave/internal/g"
	"codeberg.org/dergs/tidalwave/pkg/utils/cacheutil"
	"github.com/jwijenbergh/puregotk/v4/glib"
	"github.com/jwijenbergh/puregotk/v4/gtk"
)

type ImgUtil struct {
	cache *cacheutil.Cache
}

func (i *ImgUtil) LoadIntoImage(url string, image *gtk.Image) {
	image.Ref()
	go func() {
		loader, err := i.Load(url)
		if err != nil {
			image.Unref()
			image = nil
			return
		}

		glib.IdleAddOnce(
			g.Ptr[glib.SourceOnceFunc](func(u uintptr) {
				image.SetFromPixbuf(loader.GetAnimation().GetStaticImage())
				loader.Unref()
				image.Unref()
				image = nil
			}),
			0,
		)
	}()
}

func NewImgUtil(appId string) *ImgUtil {
	return &ImgUtil{
		cache: cacheutil.NewCache(appId, "images"),
	}
}
