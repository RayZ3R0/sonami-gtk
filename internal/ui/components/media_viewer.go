package components

import (
	"fmt"
	"time"

	"codeberg.org/dergs/tonearm/pkg/schwifty"
	. "codeberg.org/dergs/tonearm/pkg/schwifty/syntax"
	"codeberg.org/dergs/tonearm/pkg/schwifty/tracking"
	"github.com/jwijenbergh/puregotk/v4/adw"
	"github.com/jwijenbergh/puregotk/v4/gdk"
	"github.com/jwijenbergh/puregotk/v4/gio"
	"github.com/jwijenbergh/puregotk/v4/gobject"
	"github.com/jwijenbergh/puregotk/v4/gsk"
	"github.com/jwijenbergh/puregotk/v4/gtk"
)

type MediaViewer struct {
	revealer   *tracking.WeakRef // darkened background + picture container
	picture    *tracking.WeakRef
	backButton *tracking.WeakRef

	zoomLevel    float64
	baseW, baseH int

	media *gdk.Texture
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

	backButton := Button().
		ActionName("mediaviewer.close").
		IconName("go-previous-symbolic").
		ConnectConstruct(func(b *gtk.Button) {
			mvInstance.backButton = tracking.NewWeakRef(b)
		})()
	backButtonRef := tracking.NewWeakRef(backButton)

	model := gio.NewMenu()
	copyItem := gio.NewMenuItem("Copy Media", "mediaviewer.copy")
	model.AppendItem(copyItem)
	moreButton := MenuButton().
		MenuModel(&model.MenuModel).
		IconName("view-more-symbolic")
	toolbar := adw.NewToolbarView()
	toolbar.AddTopBar(HeaderBar().
		TitleWidget(Bin()).
		DecorationLayout("").
		PackStart(backButton).
		PackEnd(moreButton).
		ToGTK())
	toolbar.SetContent(
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
			Controller(&clickGesture.EventController).
			Controller(&scrollGesture.EventController).
			Controller(&zoomGesture.EventController).
			Policy(gtk.PolicyAutomaticValue, gtk.PolicyAutomaticValue).ToGTK(),
	)

	defer toolbar.Unref()

	overlay := adw.NewToastOverlay()
	overlay.SetChild(&toolbar.Widget)
	overlayRef := tracking.NewWeakRef(overlay)

	mvInstance.revealer = tracking.NewWeakRef(
		Revealer(
			VStack(overlay).
				HAlign(gtk.AlignFillValue).
				VAlign(gtk.AlignFillValue).
				HExpand(true).
				VExpand(true).
				WithCSSClass("mediaviewer"),
		).
			ConnectConstruct(func(r *gtk.Revealer) {
				actionGroup := gio.NewSimpleActionGroup()
				defer actionGroup.Unref()

				copyItem := gio.NewSimpleAction("copy", nil)
				copyItem.ConnectActivate(new(func(action gio.SimpleAction, _ uintptr) {
					mvInstance.copyMedia()

					overlayRef.Use(func(obj *gobject.Object) {
						overlay := adw.ToastOverlayNewFromInternalPtr(obj.Ptr)
						toast := adw.NewToast("Copied image to clipboard")
						toast.SetTimeout(3)
						overlay.AddToast(toast)
					})
				}))
				actionGroup.AddAction(copyItem)

				closeItem := gio.NewSimpleAction("close", nil)
				closeItem.ConnectActivate(new(func(action gio.SimpleAction, _ uintptr) {
					mvInstance.Hide()
				}))
				actionGroup.AddAction(closeItem)

				r.InsertActionGroup("mediaviewer", actionGroup)

				ShortcutController().
					ShortcutFromNames("<Ctrl>C", "mediaviewer.copy").
					ShortcutFromNames("Escape", "mediaviewer.close").
					Into(r)
			}).
			ConnectShow(func(w gtk.Widget) {
				backButtonRef.Use(func(obj *gobject.Object) {
					button := gtk.ButtonNewFromInternalPtr(obj.Ptr)
					button.GrabFocus()
				})
			}).
			Visible(false).
			TransitionType(gtk.RevealerTransitionTypeCrossfadeValue)(),
	)

	return mvInstance
}

func (mv *MediaViewer) copyMedia() {
	display := gdk.DisplayGetDefault()
	defer display.Unref()
	clipboard := display.GetClipboard()
	defer clipboard.Unref()

	clipboard.SetTexture(mv.media)
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

	snapshot := gtk.NewSnapshot()
	width := float64(paintable.GetIntrinsicWidth())
	height := float64(paintable.GetIntrinsicHeight())
	if width <= 0 {
		width = 1440
	}
	if height <= 0 {
		height = 1440
	}
	paintable.Snapshot(&snapshot.Snapshot, width, height)
	node := snapshot.FreeToNode()
	defer node.Unref()

	renderer := gsk.NewCairoRenderer()
	defer renderer.Unref()
	renderer.Realize(nil)
	defer renderer.Unrealize()
	mv.media = renderer.RenderTexture(node, nil)
}

func (mv *MediaViewer) Hide() {
	mv.revealer.Use(func(obj *gobject.Object) {
		revealer := gtk.RevealerNewFromInternalPtr(obj.Ptr)
		if !revealer.GetRevealChild() {
			return
		}

		schwifty.OnMainThreadOncePure(func() {
			mv.revealer.Use(func(obj *gobject.Object) {
				revealer := gtk.RevealerNewFromInternalPtr(obj.Ptr)
				revealer.SetRevealChild(false)
			})
		})
		time.AfterFunc(time.Duration(revealer.GetTransitionDuration()*uint(time.Millisecond)), func() {
			schwifty.OnMainThreadOncePure(func() {
				mv.revealer.Use(func(obj *gobject.Object) {
					revealer := gtk.RevealerNewFromInternalPtr(obj.Ptr)
					revealer.SetVisible(false)
				})
			})
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
