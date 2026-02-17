package lyrics

import (
	"time"

	"codeberg.org/dergs/tonearm/internal/g"
	"codeberg.org/dergs/tonearm/internal/gettext"
	"codeberg.org/dergs/tonearm/pkg/schwifty"
	"codeberg.org/dergs/tonearm/pkg/schwifty/state"
	. "codeberg.org/dergs/tonearm/pkg/schwifty/syntax"
	"github.com/jwijenbergh/puregotk/v4/graphene"
	"github.com/jwijenbergh/puregotk/v4/gtk"
)

var lyricsList = state.NewStateful[any](nil)

var lyricsView = g.Lazy(func() (w *gtk.ScrolledWindow) {
	w = ScrolledWindow().
		HPadding(16).
		BindChild(lyricsList).
		Policy(gtk.PolicyNeverValue, gtk.PolicyExternalValue)()

	adj := w.GetVadjustment()
	defer adj.Unref()
	adj.ConnectValueChanged(new(func(adj gtk.Adjustment) {
		if !scrollIsProgrammatic {
			userManuallyScrolled.SetValue(true)
		}
	}))

	return
})

var lyricsOverlay = g.Lazy(func() *gtk.Overlay {
	overlay := gtk.NewOverlay()
	overlay.SetChild(&lyricsView().Widget)
	overlay.AddOverlay(
		Button().
			HAlign(gtk.AlignEndValue).
			VAlign(gtk.AlignEndValue).
			Margin(7).
			TooltipText(gettext.Get("Sync with Track")).
			BindVisible(userManuallyScrolled).
			ConnectClicked(func(b gtk.Button) {
				if activeLyricButtonPtr := activeLyricIndex.Value(); activeLyricButtonPtr != 0 {
					w := gtk.ButtonNewFromInternalPtr(activeLyricButtonPtr)
					scrollToLyric(w)
				}
				userManuallyScrolled.SetValue(false)
			}).
			Child(
				Image().
					FromIconName("arrow-circular-top-right-symbolic"),
			).
			ToGTK(),
	)

	return overlay
})

func scrollToLyric(w *gtk.Button) {
	scrollIsProgrammatic = true
	parentWidget := w.GetParent()
	if parentWidget == nil {
		return
	}

	defer parentWidget.Unref()

	var bounds graphene.Rect
	w.ComputeBounds(parentWidget, &bounds)
	vadj := lyricsView().GetVadjustment()
	defer vadj.Unref()
	scrollViewHeight := lyricsView().GetHeight()

	// Calculate the position to center the active lyric
	widgetCenter := float64(bounds.GetY() + bounds.GetHeight()/2)
	scrollCenter := float64(scrollViewHeight / 2)
	targetPosition := widgetCenter - scrollCenter

	// Clamp the target position within valid bounds
	if targetPosition < vadj.GetLower() {
		targetPosition = vadj.GetLower()
	} else if targetPosition > vadj.GetUpper()-vadj.GetPageSize() {
		targetPosition = vadj.GetUpper() - vadj.GetPageSize()
	}

	vadj.SetValue(targetPosition)
	scrollIsProgrammatic = false
}

func setNewIndex(timing highlightTiming) {
	object := timing.Ref.Get()
	if object == nil {
		return
	}

	ptr := object.GoPointer()

	if activeLyricIndex.Value() != ptr {
		activeLyricIndex.SetValue(ptr)
	}

	if !userManuallyScrolled.Value() {
		schwifty.OnMainThreadOnce(func(uintptr) {
			w := gtk.ButtonNewFromInternalPtr(ptr)
			scrollToLyric(w)
			object.Unref()
		}, 0)
	} else {
		object.Unref()
	}
}

func setLyrics(timed bool, lyrics string, trackDuration time.Duration) {
	if timed {
		lines := parseLRCLyrics(lyrics, trackDuration)

		// Disallow user scrolling
		schwifty.OnMainThreadOncePure(func() {
			lyricsView().SetPolicy(gtk.PolicyNeverValue, gtk.PolicyExternalValue)
		})
		setLyricsLines(lines)
	} else {
		lines := parseUntimedLyrics(lyrics)

		// Allow user to scroll
		schwifty.OnMainThreadOncePure(func() {
			lyricsView().SetPolicy(gtk.PolicyNeverValue, gtk.PolicyAutomaticValue)
		})

		setLyricsLines(lines)
	}
}

func setLyricsLines(lines []any) {
	lyricsList.SetValue(
		VStack(lines...).
			Spacing(12).
			HExpand(true).
			VExpand(true),
	)
}
