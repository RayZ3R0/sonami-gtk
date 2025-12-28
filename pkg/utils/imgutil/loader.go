package imgutil

import (
	"github.com/jwijenbergh/puregotk/v4/gdkpixbuf"
)

func (i *ImgUtil) Load(url string) (*gdkpixbuf.PixbufLoader, error) {
	image, err := i.fetch(url)
	if err != nil {
		return nil, err
	}

	loader := gdkpixbuf.NewPixbufLoader()
	if _, err := loader.Write(image, uint(len(image))); err != nil {
		loader.Close()
		return nil, err
	}

	if _, err := loader.Close(); err != nil {
		return nil, err
	}

	return loader, nil
}

func (i *ImgUtil) LoadPixbuf(url string) (*gdkpixbuf.Pixbuf, error) {
	loader, err := i.Load(url)
	if err != nil {
		return nil, err
	}
	defer loader.Unref()

	return loader.GetAnimation().GetStaticImage(), nil
}
