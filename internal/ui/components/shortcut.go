package components

import (
	"context"

	. "codeberg.org/dergs/tidalwave/pkg/gui"
	"github.com/diamondburned/gotk4/pkg/gtk/v4"
	"github.com/diamondburned/gotkit/gtkutil/imgutil"
	"github.com/infinytum/injector"
)

type Shortcut struct {
	*gtk.Button
	cover    *gtk.Image
	subTitle *TextImpl
	title    *TextImpl
}

func (s *Shortcut) LoadCover(url string) {
	imgutil.AsyncGET(injector.MustInject[context.Context](), url, imgutil.ImageSetterFromImage(s.cover))
}

func (s *Shortcut) SetSubTitle(subtitle string) *Shortcut {
	s.subTitle.Text(subtitle)
	if subtitle == "" {
		s.subTitle.GTKWidget().SetVisible(false)
	} else {
		s.subTitle.GTKWidget().SetVisible(true)
	}
	return s
}

func (s *Shortcut) SetTitle(title string) *Shortcut {
	s.title.Text(title)
	return s
}

func NewShortcut() *Shortcut {
	shortcut := &Shortcut{
		cover:    gtk.NewImageFromResource("/org/codeberg/dergs/tidalwave/icons/scalable/state/missing-album.svg"),
		subTitle: Text("").HAlign(gtk.AlignStart).FontWeight(300),
		title:    Text("").HAlign(gtk.AlignStart),
	}
	shortcut.cover.SetPixelSize(54)
	shortcut.cover.SetHExpand(true)
	shortcut.cover.SetHAlign(gtk.AlignEnd)
	shortcut.cover.SetOverflow(gtk.OverflowHidden)
	shortcut.subTitle.GTKWidget().SetVisible(false)

	shortcut.Button = gtk.NewButton()
	shortcut.Button.SetHAlign(gtk.AlignStart)
	shortcut.Button.SetHExpand(false)
	shortcut.Button.SetChild(
		HStack(
			VStack(
				shortcut.title,
				shortcut.subTitle,
			).HAlign(gtk.AlignStart).VAlign(gtk.AlignCenter),
			Wrapper(shortcut.cover).CornerRadius(10),
		).CSS("box { min-width: 400px; }"),
	)
	Wrapper(shortcut).CSS("button { padding-left: 15px; padding-right: 5px; padding-top: 5px; padding-bottom: 5px;}")

	return shortcut
}
