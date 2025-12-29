package horizontal_list

import (
	"math"

	"codeberg.org/dergs/tidalwave/pkg/schwifty"
	"codeberg.org/dergs/tidalwave/pkg/schwifty/state"
	. "codeberg.org/dergs/tidalwave/pkg/schwifty/syntax"
	"github.com/jwijenbergh/puregotk/v4/gtk"
)

type HorizontalList struct {
	schwifty.Box

	container   *gtk.Box
	marginState *state.State[int]
}

func (h *HorizontalList) Append(child any) *HorizontalList {
	h.container.Append(schwifty.ResolveWidget(child))
	return h
}

func (h *HorizontalList) SetPageMargin(margin int) *HorizontalList {
	h.marginState.SetValue(margin)
	return h
}

func NewHorizontalList(title string) *HorizontalList {
	marginState := state.NewStateful[int](0)
	container := HStack().BindHMargin(marginState)()

	var hAdjust *gtk.Adjustment
	nextButton := Button().
		MinHeight(10).MinWidth(10).HPadding(10).
		Child(Image().FromIconName("go-next-symbolic").PixelSize(10)).
		ConnectClicked(func(b gtk.Button) {
			current := hAdjust.GetValue()
			current -= math.Mod(current, 192)
			hAdjust.SetStepIncrement(192)
			hAdjust.SetValue(current + hAdjust.GetStepIncrement())
		})

	previousButton := Button().
		MinHeight(10).MinWidth(10).HPadding(10).
		Child(Image().FromIconName("go-previous-symbolic").PixelSize(10)).
		ConnectClicked(func(b gtk.Button) {
			current := hAdjust.GetValue()
			if math.Mod(current, 192) > 0 {
				current += 192 - math.Mod(current, 192)
			}
			hAdjust.SetStepIncrement(192)
			hAdjust.SetValue(current - hAdjust.GetStepIncrement())
		})

	return &HorizontalList{
		Box: VStack(
			HStack(
				HStack(
					Label(title).
						FontWeight(600).
						FontSize(20).
						VAlign(gtk.AlignCenterValue),
					Spacer().VExpand(false),
					HStack(
						previousButton,
						nextButton,
					).Spacing(10),
				).BindHMargin(marginState),
			).HMargin(10).MarginBottom(10),
			ScrolledWindow().
				Child(ManagedWidget(&container.Widget)).
				VAlign(gtk.AlignStartValue).
				Policy(gtk.PolicyExternalValue, gtk.PolicyNeverValue).
				PropagateNaturalWidth(true).
				PropagateNaturalWidth(true).
				ConnectConstruct(func(sw *gtk.ScrolledWindow) {
					hAdjust = sw.GetHadjustment()
				}),
		),
		marginState: marginState,
		container:   container,
	}
}
