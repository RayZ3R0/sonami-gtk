package ui

import (
	"unsafe"

	"codeberg.org/dergs/tidalwave/internal/g"
	"codeberg.org/dergs/tidalwave/internal/player"
	"codeberg.org/dergs/tidalwave/internal/router"
	"github.com/jwijenbergh/puregotk/v4/gio"
	"github.com/jwijenbergh/puregotk/v4/glib"
)

func (w *Window) installActions() {
	navigateBackAction := gio.NewSimpleAction("navigate-back", nil)
	navigateBackAction.ConnectActivate(g.Ptr(func(action gio.SimpleAction, parameter uintptr) {
		router.Back()
	}))
	w.AddAction(navigateBackAction)
	w.GetApplication().SetAccelsForAction("win.navigate-back", []string{"<Alt>Left"})

	playTrackAction := gio.NewSimpleAction("player.play-track", glib.NewVariantType("s"))
	playTrackAction.ConnectActivate(g.Ptr(func(action gio.SimpleAction, parameter uintptr) {
		variant := (*glib.Variant)(unsafe.Pointer(parameter))
		id := variant.GetString(nil)
		player.Play(id)
	}))
	w.AddAction(playTrackAction)

	routeAlbumAction := gio.NewSimpleAction("route.album", glib.NewVariantType("s"))
	routeAlbumAction.ConnectActivate(g.Ptr(func(action gio.SimpleAction, parameter uintptr) {
		variant := (*glib.Variant)(unsafe.Pointer(parameter))
		router.Navigate("album", router.Params{
			"id": variant.GetString(nil),
		})
	}))
	w.AddAction(routeAlbumAction)

	routePlaylistAction := gio.NewSimpleAction("route.playlist", glib.NewVariantType("s"))
	routePlaylistAction.ConnectActivate(g.Ptr(func(action gio.SimpleAction, parameter uintptr) {
		variant := (*glib.Variant)(unsafe.Pointer(parameter))
		router.Navigate("playlist", router.Params{
			"id": variant.GetString(nil),
		})
	}))
	w.AddAction(routePlaylistAction)

	routeArtistAction := gio.NewSimpleAction("route.artist", glib.NewVariantType("s"))
	routeArtistAction.ConnectActivate(g.Ptr(func(action gio.SimpleAction, parameter uintptr) {
		variant := (*glib.Variant)(unsafe.Pointer(parameter))
		router.Navigate("artist", router.Params{
			"id": variant.GetString(nil),
		})
	}))
	w.AddAction(routeArtistAction)
}
