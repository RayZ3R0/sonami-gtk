package ui

import (
	"context"
	"unsafe"

	"codeberg.org/dergs/tonearm/internal/g"
	"codeberg.org/dergs/tonearm/internal/notifications"
	"codeberg.org/dergs/tonearm/internal/player"
	"codeberg.org/dergs/tonearm/internal/router"
	"codeberg.org/dergs/tonearm/internal/secrets"
	"codeberg.org/dergs/tonearm/internal/settings"
	"codeberg.org/dergs/tonearm/internal/ui/components/linking"
	"codeberg.org/dergs/tonearm/pkg/schwifty"
	"codeberg.org/dergs/tonearm/pkg/tidalapi"
	"codeberg.org/dergs/tonearm/pkg/tidalapi/auth"
	"github.com/jwijenbergh/puregotk/v4/adw"
	"github.com/jwijenbergh/puregotk/v4/gio"
	"github.com/jwijenbergh/puregotk/v4/glib"
)

func (w *Window) installActions() {
	aboutAction := gio.NewSimpleAction("about", nil)
	aboutAction.ConnectActivate(g.Ptr(func(action gio.SimpleAction, parameter uintptr) {
		w.PresentAbout()
	}))
	w.GetApplication().Application.AddAction(aboutAction)

	preferencesAction := gio.NewSimpleAction("preferences", nil)
	preferencesAction.ConnectActivate(g.Ptr(func(action gio.SimpleAction, parameter uintptr) {
		w.PresentPreferences()
	}))
	w.GetApplication().Application.AddAction(preferencesAction)
	w.GetApplication().SetAccelsForAction("app.preferences", []string{"<Control>comma"})

	quitAction := gio.NewSimpleAction("quit", nil)
	quitAction.ConnectActivate(g.Ptr(func(action gio.SimpleAction, parameter uintptr) {
		w.GetApplication().Quit()
	}))
	w.GetApplication().Application.AddAction(quitAction)
	w.GetApplication().SetAccelsForAction("app.quit", []string{"<Ctrl>q"})

	navigateBackAction := gio.NewSimpleAction("navigate-back", nil)
	navigateBackAction.ConnectActivate(g.Ptr(func(action gio.SimpleAction, parameter uintptr) {
		router.Back()
	}))
	w.AddAction(navigateBackAction)
	w.GetApplication().SetAccelsForAction("win.navigate-back", []string{"<Alt>Left"})

	searchAction := gio.NewSimpleAction("search", nil)
	searchAction.ConnectActivate(g.Ptr(func(action gio.SimpleAction, parameter uintptr) {
		router.Navigate("search")
	}))
	w.AddAction(searchAction)
	w.GetApplication().SetAccelsForAction("win.search", []string{"<Ctrl>f"})

	playTrackAction := gio.NewSimpleAction("player.play-track", glib.NewVariantType("s"))
	playTrackAction.ConnectActivate(g.Ptr(func(action gio.SimpleAction, parameter uintptr) {
		variant := (*glib.Variant)(unsafe.Pointer(parameter))
		id := variant.GetString(nil)
		go player.PlayTrack(id)
	}))
	w.AddAction(playTrackAction)

	playPlaylistAction := gio.NewSimpleAction("player.play-playlist", glib.NewVariantType("s"))
	playPlaylistAction.ConnectActivate(g.Ptr(func(action gio.SimpleAction, parameter uintptr) {
		variant := (*glib.Variant)(unsafe.Pointer(parameter))
		id := variant.GetString(nil)
		go player.PlayPlaylist(id, false, "")
	}))
	w.AddAction(playPlaylistAction)

	shuffleAction := gio.NewSimpleAction("player.shuffle", nil)
	shuffleAction.ConnectActivate(g.Ptr(func(action gio.SimpleAction, parameter uintptr) {
		go player.ToggleShuffle()
	}))
	w.AddAction(shuffleAction)

	nextAction := gio.NewSimpleAction("player.next", nil)
	nextAction.ConnectActivate(g.Ptr(func(action gio.SimpleAction, parameter uintptr) {
		go player.Next()
	}))
	w.AddAction(nextAction)

	previousAction := gio.NewSimpleAction("player.previous", nil)
	previousAction.ConnectActivate(g.Ptr(func(action gio.SimpleAction, parameter uintptr) {
		go player.Previous()
	}))
	w.AddAction(previousAction)

	repeatAction := gio.NewSimpleAction("player.repeat", nil)
	repeatAction.ConnectActivate(g.Ptr(func(action gio.SimpleAction, parameter uintptr) {
		player.CycleRepeatMode()
	}))
	w.AddAction(repeatAction)

	queueTrackAction := gio.NewSimpleAction("player.queue-track", glib.NewVariantType("s"))
	queueTrackAction.ConnectActivate(g.Ptr(func(action gio.SimpleAction, parameter uintptr) {
		variant := (*glib.Variant)(unsafe.Pointer(parameter))
		id := variant.GetString(nil)
		go player.UserQueue.AddTrackID(id, false)
	}))
	w.AddAction(queueTrackAction)

	routeAlbumAction := gio.NewSimpleAction("route.album", glib.NewVariantType("s"))
	routeAlbumAction.ConnectActivate(g.Ptr(func(action gio.SimpleAction, parameter uintptr) {
		variant := (*glib.Variant)(unsafe.Pointer(parameter))
		router.Navigate("album/" + variant.GetString(nil))
	}))
	w.AddAction(routeAlbumAction)

	routePlaylistAction := gio.NewSimpleAction("route.playlist", glib.NewVariantType("s"))
	routePlaylistAction.ConnectActivate(g.Ptr(func(action gio.SimpleAction, parameter uintptr) {
		variant := (*glib.Variant)(unsafe.Pointer(parameter))
		router.Navigate("playlist/" + variant.GetString(nil))
	}))
	w.AddAction(routePlaylistAction)

	routeArtistAction := gio.NewSimpleAction("route.artist", glib.NewVariantType("s"))
	routeArtistAction.ConnectActivate(g.Ptr(func(action gio.SimpleAction, parameter uintptr) {
		variant := (*glib.Variant)(unsafe.Pointer(parameter))
		router.Navigate("artist/" + variant.GetString(nil))
	}))
	w.AddAction(routeArtistAction)

	signInAction := gio.NewSimpleAction("sign-in", nil)
	signInAction.ConnectActivate(g.Ptr(func(action gio.SimpleAction, parameter uintptr) {
		go func() {
			var dialog *adw.AlertDialog
			resp, err := tidalapi.StartDeviceLinking(func(dlc *auth.DeviceLinkingChallenge, cancel context.CancelFunc) {
				schwifty.OnMainThreadOnce(func(u uintptr) {
					dialog = linking.NewLinking(&w.Window, dlc.UserCode, dlc.VerificationUriComplete, cancel)
					defer dialog.Unref()
					dialog.Present(&w.Widget)
				}, 0)
			})
			defer dialog.ForceClose()
			if err != nil {
				notifications.OnToast.Notify("Sign in failed or aborted")
				return
			}
			secrets.SetRefreshToken(resp.RefreshToken)
			notifications.OnToast.Notify("Signed in as " + resp.User.Email)
			router.Refresh()
		}()
	}))
	w.AddAction(signInAction)

	signOutAction := gio.NewSimpleAction("sign-out", nil)
	signOutAction.ConnectActivate(g.Ptr(func(action gio.SimpleAction, parameter uintptr) {
		go func() {
			secrets.DeleteRefreshToken()
			notifications.OnToast.Notify("Signed out")
			router.Refresh()
		}()
	}))
	w.AddAction(signOutAction)

	setAsDefaultAction := gio.NewSimpleAction("set-as-default", nil)
	setAsDefaultAction.ConnectActivate(g.Ptr(func(action gio.SimpleAction, parameter uintptr) {
		settings.General().SetDefaultPage(router.Current().Path)
		notifications.OnToast.Notify("Set current page as default")
	}))
	w.AddAction(setAsDefaultAction)
}
