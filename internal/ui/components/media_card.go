package components

import (
	"context"
	"fmt"
	"math/rand/v2"

	"github.com/diamondburned/gotk4/pkg/gtk/v4"
	"github.com/diamondburned/gotk4/pkg/pango"
	"github.com/diamondburned/gotkit/gtkutil/cssutil"
	"github.com/diamondburned/gotkit/gtkutil/imgutil"
	"github.com/infinytum/injector"
)

var mediaCardCSS = cssutil.Applier("media-card", `
	.media-card {
		padding: 10px;
	}
	.media-card:not(:hover) {
		background-color: transparent;
	}

	.picture-frame {
		border-radius: 12px;
	}

	.media-card-title {
		font-size: 16px;
	}

	.media-card-subtitle {
		font-weight: 400;
		font-size: 14px;
		color: #939393;
	}
`)

type MediaCard struct {
	*gtk.Button
	image         *gtk.Image
	subTitleLabel *gtk.Label
	titleLabel    *gtk.Label
}

func (m *MediaCard) LoadImage(url string) {
	imgutil.AsyncGET(injector.MustInject[context.Context](), url, imgutil.ImageSetterFromImage(m.image))
}

func (m *MediaCard) SetTitle(title string) {
	m.titleLabel.SetText(title)
	m.titleLabel.SetTooltipText(title)
}

func (m *MediaCard) SetSubTitle(artist string) {
	m.subTitleLabel.SetText(artist)
	m.subTitleLabel.SetTooltipText(artist)
}

func NewMediaCard() *MediaCard {
	cover := gtk.NewImage()
	cover.SetPixelSize(172)

	frame := gtk.NewAspectFrame(0.5, 0.5, 1, false)
	frame.SetChild(cover)
	frame.AddCSSClass("picture-frame")
	frame.SetOverflow(gtk.OverflowHidden)

	titleLabel := gtk.NewLabel("Album " + fmt.Sprintf("%d", rand.IntN(100)))
	titleLabel.SetHAlign(gtk.AlignStart)
	titleLabel.AddCSSClass("media-card-title")
	titleLabel.SetMarginTop(10)
	titleLabel.SetEllipsize(pango.EllipsizeEnd)

	subTitleLabel := gtk.NewLabel("Embark Studios\n2024")
	subTitleLabel.SetHAlign(gtk.AlignStart)
	subTitleLabel.AddCSSClass("media-card-subtitle")
	subTitleLabel.SetEllipsize(pango.EllipsizeEnd)
	subTitleLabel.SetMarginTop(2)

	box := gtk.NewBox(gtk.OrientationVertical, 0)
	box.Append(frame)
	box.Append(titleLabel)
	box.Append(subTitleLabel)
	box.SetHAlign(gtk.AlignStart)

	button := gtk.NewButton()
	button.AddCSSClass("image-button")
	button.SetChild(box)
	button.SetHExpand(false)
	button.SetVExpand(false)

	mediaCard := &MediaCard{
		button,
		cover,
		subTitleLabel,
		titleLabel,
	}

	mediaCardCSS(mediaCard)

	return mediaCard
}
