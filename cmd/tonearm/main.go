package main

import (
	"log/slog"
	"os"
	"strings"

	_ "codeberg.org/dergs/tonearm/internal/log"

	_ "codeberg.org/dergs/tonearm/internal/icons"
	_ "codeberg.org/dergs/tonearm/internal/styles"

	// _ "codeberg.org/dergs/tonearm/internal/features/scrobbling"
	_ "codeberg.org/dergs/tonearm/internal/services"

	"codeberg.org/dergs/tonearm/internal/g"
	"codeberg.org/dergs/tonearm/internal/player"
	"codeberg.org/dergs/tonearm/internal/router"
	"codeberg.org/dergs/tonearm/internal/secrets"
	"codeberg.org/dergs/tonearm/internal/settings"
	"codeberg.org/dergs/tonearm/internal/ui"
	"codeberg.org/dergs/tonearm/pkg/mpris"
	"codeberg.org/dergs/tonearm/pkg/schwifty/tracking"
	"codeberg.org/dergs/tonearm/pkg/utils/imgutil"
	"github.com/infinytum/injector"
	"github.com/jwijenbergh/puregotk/v4/adw"
	"github.com/jwijenbergh/puregotk/v4/gio"
	"github.com/jwijenbergh/puregotk/v4/gtk"
)

func init() {
	slog.SetLogLoggerLevel(slog.LevelInfo)
	if os.Getenv("TONEARM_DEBUG") == "1" {
		go tracking.LogAliveObjects()
	}
}

var app *adw.Application

func main() {
	app = adw.NewApplication("dev.dergs.Tonearm", gio.GApplicationHandlesCommandLineValue)
	defer app.Unref()
	app.ConnectActivate(g.Ptr(onActivate))
	app.ConnectCommandLine(g.Ptr(onCommandLine))
	injector.Singleton(func() *adw.Application {
		return app
	})

	injector.Singleton(func(app *adw.Application) *imgutil.ImgUtil {
		return imgutil.NewImgUtil(app.GetApplicationId())
	})

	injector.DeferredSingleton(func() *mpris.Server {
		mprisServer := mpris.NewMprisServer("org.mpris.MediaPlayer2."+app.GetApplicationId(), app.GetApplicationId(), "Tonearm")
		mprisServer.OnPause(player.Pause)
		mprisServer.OnPlayPause(player.PlayPause)
		mprisServer.OnPlay(player.Play)
		mprisServer.OnTrackNext(player.Next)
		mprisServer.OnTrackPrevious(player.Previous)
		mprisServer.OnQuit(func() { quit(0) })
		mprisServer.OnRaise(func() {
			window := injector.MustInject[*adw.ApplicationWindow]()
			window.Show()
			window.Present()
		})
		mprisServer.OnLoopStatusChanged(func(loopStatus mpris.LoopStatus) {
			switch loopStatus {
			case mpris.LoopNone:
				go player.SetRepeatMode(player.RepeatModeNone)
			case mpris.LoopTrack:
				go player.SetRepeatMode(player.RepeatModeTrack)
			case mpris.LoopPlaylist:
				go player.SetRepeatMode(player.RepeatModeQueue)
			}
		})
		mprisServer.OnSeek(player.SeekToPositionRelative)
		mprisServer.OnSetPosition(player.SeekToPosition)
		mprisServer.OnVolumeChanged(func(newVal float64) {
			player.SetVolume(newVal)
		})
		mprisServer.OnShuffleChanged(func(shuffle bool) {
			player.SetShuffle(shuffle)
		})
		mprisServer.Export()
		return mprisServer
	})

	if code := app.Run(len(os.Args), os.Args); code > 0 {
		quit(code)
	}
}

func quit(code int) {
	app.Quit()
	os.Exit(code)
}

func onActivate(_ gio.Application) {
	window := ui.NewWindow(app)
	window.Present()
	window.ConnectCloseRequest(g.Ptr(func(w gtk.Window) bool {
		// Only allow running in background if there is a track playing,
		// so that the user can bring the app back up with MPRIS
		if settings.General().ShouldRunInBackground() && player.TrackChanged.CurrentValue() != nil {
			window.Hide()
			return true
		}
		return false
	}))

	if err := secrets.Healthcheck(); err != nil {
		slog.Error("Secret service health check failed", "title", err.Title, "body", err.Body, "fatal", err.Fatal)
		window.PresentSecretServiceError(err)
	}

}

var isActive bool

func onCommandLine(app gio.Application, ptr uintptr) int {
	cmd := gio.ApplicationCommandLineNewFromInternalPtr(ptr)
	args := cmd.GetArguments(nil)
	if !isActive {
		app.Activate()
		isActive = true
	} else {
		window := injector.MustInject[*adw.ApplicationWindow]()
		window.Show()
	}

	if len(args) == 2 {
		url := args[1]
		if strings.HasPrefix(url, "tidal://track/") {
			player.PlayTrack(strings.TrimPrefix(url, "tidal://track/"))
		} else {
			router.Navigate(url)
		}
	}
	return 0
}
