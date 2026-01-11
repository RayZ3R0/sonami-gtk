package main

import (
	"log/slog"
	"os"
	"strings"

	_ "codeberg.org/dergs/tidalwave/internal/features/scrobbling"
	"codeberg.org/dergs/tidalwave/internal/g"
	_ "codeberg.org/dergs/tidalwave/internal/icons"
	"codeberg.org/dergs/tidalwave/internal/player"
	"codeberg.org/dergs/tidalwave/internal/router"
	"codeberg.org/dergs/tidalwave/internal/secrets"
	"codeberg.org/dergs/tidalwave/internal/settings"
	_ "codeberg.org/dergs/tidalwave/internal/styles"
	"codeberg.org/dergs/tidalwave/internal/ui"
	"codeberg.org/dergs/tidalwave/pkg/mpris"
	"codeberg.org/dergs/tidalwave/pkg/schwifty/tracking"
	"codeberg.org/dergs/tidalwave/pkg/tidalapi"
	"codeberg.org/dergs/tidalwave/pkg/utils/imgutil"
	"github.com/infinytum/injector"
	"github.com/jwijenbergh/puregotk/v4/adw"
	"github.com/jwijenbergh/puregotk/v4/gio"
	"github.com/jwijenbergh/puregotk/v4/gtk"
)

func init() {
	slog.SetLogLoggerLevel(slog.LevelInfo)
	if os.Getenv("TIDAL_WAVE_DEBUG") == "1" {
		slog.SetLogLoggerLevel(slog.LevelDebug)
		go tracking.LogAliveWidgets()
	}
}

var app *adw.Application

func main() {
	app = adw.NewApplication("org.codeberg.dergs.tidalwave", gio.GApplicationHandlesCommandLineValue)
	defer app.Unref()
	app.ConnectActivate(g.Ptr(onActivate))
	app.ConnectCommandLine(g.Ptr(onCommandLine))
	injector.Singleton(func() *adw.Application {
		return app
	})

	injector.Singleton(func() *tidalapi.TidalAPI {
		countryCode, err := tidalapi.FetchCountryCode()
		if err != nil {
			slog.Error("Failed to fetch country code, defaulting to WW", err)
			countryCode = "WW"
		}
		slog.Info("Discovered country code", "countryCode", countryCode)
		return tidalapi.NewClient(countryCode, secrets.NewTokenAuthStrategy())
	})

	injector.Singleton(func(app *adw.Application) *imgutil.ImgUtil {
		return imgutil.NewImgUtil(app.GetApplicationId())
	})

	injector.DeferredSingleton(func() *mpris.Server {
		mprisServer := mpris.NewMprisServer("org.mpris.MediaPlayer2.TidalWave", app.GetApplicationId(), "Tidal Wave")
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
