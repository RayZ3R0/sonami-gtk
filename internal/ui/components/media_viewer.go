package components

import (
	"fmt"
	"time"

	"codeberg.org/dergs/tonearm/pkg/schwifty"
	. "codeberg.org/dergs/tonearm/pkg/schwifty/syntax"
	"codeberg.org/dergs/tonearm/pkg/schwifty/tracking"
	"github.com/jwijenbergh/puregotk/v4/gdk"
	"github.com/jwijenbergh/puregotk/v4/gobject"
	"github.com/jwijenbergh/puregotk/v4/gtk"
)

type MediaViewer struct {
	revealer     *tracking.WeakRef // darkened background + picture container
	picture      *tracking.WeakRef
	zoomLevel    float64
	baseW, baseH int
}

var (
	mvInstance *MediaViewer
)

func (mv *MediaViewer) ToGTK() *gtk.Widget {
	obj := mv.revealer.Get()
	revealer := gtk.RevealerNewFromInternalPtr(obj.Ptr)
	return &revealer.Widget
}

func GetMediaViewer() *MediaViewer {
	if mvInstance != nil {
		return mvInstance
	}

	mvInstance = &MediaViewer{}

	clickGesture := gtk.NewGestureClick()
	clickGesture.ConnectPressed(new(func(_ gtk.GestureClick, _ int, _ float64, _ float64) {
		mvInstance.Hide()
	}))

	scrollGesture := gtk.NewEventControllerScroll(gtk.EventControllerScrollVerticalValue)
	scrollGesture.SetPropagationPhase(gtk.PhaseCaptureValue)
	scrollGesture.ConnectScroll(new(func(_ gtk.EventControllerScroll, _, dy float64) bool {
		// Only zoom when Ctrl is held
		state := scrollGesture.GetCurrentEventState()
		if state&gdk.ControlMaskValue != 0 {
			if dy < 0 {
				go mvInstance.zoom(1.1)
			} else {
				go mvInstance.zoom(0.9)
			}
			return true // consume the event
		}
		return false
	}))

	zoomGesture := gtk.NewGestureZoom()
	var startZoom float64
	zoomGesture.ConnectBegin(new(func(_ gtk.Gesture, _ uintptr) {
		startZoom = mvInstance.zoomLevel
	}))
	zoomGesture.ConnectScaleChanged(new(func(_ gtk.GestureZoom, scale float64) {
		go mvInstance.setZoom(startZoom * scale)
	}))

	// --- Double-click on picture to reset zoom ---
	dblClick := gtk.NewGestureClick()
	dblClick.SetButton(1)
	dblClick.ConnectPressed(new(func(gesture gtk.GestureClick, nPress int, _, _ float64) {
		gesture.SetState(gtk.EventSequenceClaimedValue)
		if nPress == 2 {
			mvInstance.resetZoom()
		}
	}))

	mvInstance.revealer = tracking.NewWeakRef(
		Revealer(
			VStack(
				ScrolledWindow().
					Child(
						Picture().
							HAlign(gtk.AlignCenterValue).
							VAlign(gtk.AlignCenterValue).
							CanShrink(true).
							ContentFit(gtk.ContentFitContainValue).
							ConnectConstruct(func(p *gtk.Picture) {
								mvInstance.picture = tracking.NewWeakRef(p)
							}).
							Controller(&dblClick.EventController),
					).
					HExpand(true).
					VExpand(true).
					Policy(gtk.PolicyAutomaticValue, gtk.PolicyAutomaticValue),
			).
				HAlign(gtk.AlignFillValue).
				VAlign(gtk.AlignFillValue).
				HExpand(true).
				VExpand(true).
				Controller(&clickGesture.EventController).
				Controller(&scrollGesture.EventController).
				Controller(&zoomGesture.EventController).
				WithCSSClass("mediaviewer"),
		).
			Visible(false).
			TransitionType(gtk.RevealerTransitionTypeCrossfadeValue)(),
	)

	return mvInstance
}

func (mv *MediaViewer) ShowFile(paintable gdk.Paintable) {
	if paintable == nil {
		return
	}

	mv.resetZoom()
	mv.picture.Use(func(obj *gobject.Object) {
		picture := gtk.PictureNewFromInternalPtr(obj.Ptr)
		picture.SetPaintable(paintable)
	})

	mv.revealer.Use(func(obj *gobject.Object) {
		revealer := gtk.RevealerNewFromInternalPtr(obj.Ptr)
		revealer.SetVisible(true)
		revealer.SetRevealChild(true)
	})
}

func (mv *MediaViewer) Hide() {
	mv.revealer.Use(func(obj *gobject.Object) {
		revealer := gtk.RevealerNewFromInternalPtr(obj.Ptr)
		if !revealer.GetRevealChild() {
			return
		}

		revealer.SetRevealChild(false)
		time.AfterFunc(time.Duration(revealer.GetTransitionDuration()*uint(time.Millisecond)), func() {
			revealer.SetVisible(false)
		})
	})
}

func (mv *MediaViewer) resetZoom() {
	mv.zoomLevel = 1.0
	mv.baseW = 0 // force re-fetch of intrinsic size on next zoom
	mv.baseH = 0

	schwifty.OnMainThreadOncePure(func() {
		mv.picture.Use(func(obj *gobject.Object) {
			picture := gtk.PictureNewFromInternalPtr(obj.Ptr)
			picture.SetSizeRequest(-1, -1)
		})
	})
}

func (mv *MediaViewer) zoom(factor float64) {
	mv.setZoom(mv.zoomLevel * factor)
}

func (mv *MediaViewer) setZoom(level float64) {
	// Clamp zoom between 0.5x and 5x
	if level < 1 {
		level = 1
	}
	if level > 10.0 {
		level = 10.0
	}
	mv.zoomLevel = level

	if mv.baseW == 0 || mv.baseH == 0 {
		c := make(chan struct{})
		schwifty.OnMainThreadOncePure(func() {
			mv.picture.Use(func(obj *gobject.Object) {
				pic := gtk.PictureNewFromInternalPtr(obj.Ptr)
				mv.baseW = pic.GetSize(gtk.OrientationHorizontalValue)
				mv.baseH = pic.GetSize(gtk.OrientationVerticalValue)
			})

			close(c)
		})

		<-c
	}

	w := int(float64(mv.baseW) * mv.zoomLevel)
	h := int(float64(mv.baseH) * mv.zoomLevel)
	schwifty.OnMainThreadOncePure(func() {
		mv.picture.Use(func(obj *gobject.Object) {
			fmt.Println(w, h)
			pic := gtk.PictureNewFromInternalPtr(obj.Ptr)
			pic.SetSizeRequest(w, h)
		})
	})
}
