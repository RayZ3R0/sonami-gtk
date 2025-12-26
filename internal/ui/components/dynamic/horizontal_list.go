package dynamic

import (
	"math"

	"codeberg.org/dergs/tidalwave/pkg/gui"
	. "codeberg.org/dergs/tidalwave/pkg/gui"
	"github.com/diamondburned/gotk4/pkg/gtk/v4"
	"github.com/diamondburned/gotkit/gtkutil/cssutil"
)

var horizontalListButtonsCSS = cssutil.Applier("media-carousel-button", `
	.media-carousel-button {
		min-height: 10px;
		min-width: 10px;
	  	padding-left: 10px;
		padding-right: 10px;
	}
`)

type HorizontalList struct {
	*BoxImpl
	container  *gui.BoxImpl
	titleLabel *gui.TextImpl
	topBar     *gui.BoxImpl
}

func (h *HorizontalList) Append(child gtk.Widgetter) *HorizontalList {
	h.container.Append(child)
	return h
}

func (h *HorizontalList) HMargin(margin int) *HorizontalList {
	h.container.HMargin(margin)
	h.topBar.HMargin(margin)
	return h
}

func (h *HorizontalList) SetTitle(title string) *HorizontalList {
	h.titleLabel.Text(title)
	return h
}

func NewHorizontalList() *HorizontalList {
	titleLabel := Text("List")
	container := HStack().HMargin(40)

	scrolledWindow := gtk.NewScrolledWindow()
	scrolledWindow.SetChild(container)
	scrolledWindow.SetVAlign(gtk.AlignStart)
	scrolledWindow.SetPropagateNaturalHeight(true)
	scrolledWindow.SetPropagateNaturalWidth(true)
	scrolledWindow.SetPolicy(gtk.PolicyExternal, gtk.PolicyNever)
	hadj := scrolledWindow.HAdjustment()

	nextIcon := gtk.NewImageFromIconName("go-next-symbolic")
	nextIcon.SetPixelSize(10)
	rightButton := gtk.NewButton()
	rightButton.SetChild(nextIcon)
	rightButton.SetHAlign(gtk.AlignEnd)
	rightButton.AddCSSClass("media-carousel-button")
	rightButton.ConnectClicked(func() {
		current := hadj.Value()
		current -= math.Mod(current, 192)
		hadj.SetStepIncrement(192)
		hadj.SetValue(current + hadj.StepIncrement())
	})
	horizontalListButtonsCSS(rightButton)

	backIcon := gtk.NewImageFromIconName("go-previous-symbolic")
	backIcon.SetPixelSize(10)
	leftButton := gtk.NewButton()
	leftButton.SetChild(backIcon)
	leftButton.SetHAlign(gtk.AlignEnd)
	leftButton.AddCSSClass("media-carousel-button")
	leftButton.ConnectClicked(func() {
		current := hadj.Value()
		if math.Mod(current, 192) > 0 {
			current += 192 - math.Mod(current, 192)
		}
		hadj.SetStepIncrement(192)
		hadj.SetValue(current - hadj.StepIncrement())
	})
	horizontalListButtonsCSS(leftButton)

	topBar := HStack(
		titleLabel.
			FontWeight(600).
			FontSize(20).
			VAlign(gtk.AlignCenter),
		Spacer().VExpand(false),
		HStack(
			leftButton,
			rightButton,
		).Spacing(10),
	)

	h := &HorizontalList{
		VStack(
			HStack(topBar).MarginLeft(10).MarginBottom(10).MarginRight(10),
			scrolledWindow,
		),
		container,
		titleLabel,
		topBar,
	}

	return h
}
