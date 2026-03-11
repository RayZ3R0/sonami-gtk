package main

import (
	"log/slog"
	"os"

	_ "github.com/RayZ3R0/sonami-gtk/internal/log"

	_ "github.com/RayZ3R0/sonami-gtk/internal/icons"
	_ "github.com/RayZ3R0/sonami-gtk/internal/styles"

	_ "github.com/RayZ3R0/sonami-gtk/internal/features/discord"
	_ "github.com/RayZ3R0/sonami-gtk/internal/features/scrobbling"
	_ "github.com/RayZ3R0/sonami-gtk/internal/services"

	"codeberg.org/puregotk/puregotk/v4/adw"
	"codeberg.org/puregotk/puregotk/v4/gio"
	"codeberg.org/puregotk/puregotk/v4/gtk"
	"github.com/RayZ3R0/sonami-gtk/internal/player"
	"github.com/RayZ3R0/sonami-gtk/internal/router"
	"github.com/RayZ3R0/sonami-gtk/internal/settings"
	"github.com/RayZ3R0/sonami-gtk/internal/ui"
	"github.com/RayZ3R0/sonami-gtk/pkg/schwifty/tracking"
	"github.com/RayZ3R0/sonami-gtk/pkg/utils/imgutil"
	"github.com/infinytum/injector"
)

func init() {
	slog.SetLogLoggerLevel(slog.LevelInfo)
	if os.Getenv("SONAMI_DEBUG") == "1" {
		go tracking.LogAliveObjects()
	}
}

var app *adw.Application

func main() {
	app = adw.NewApplication("io.github.rayz3r0.SonamiGtk", gio.GApplicationHandlesCommandLineValue)
	defer app.Unref()
	app.ConnectActivate(new(onActivate))
	app.ConnectCommandLine(new(onCommandLine))
	injector.Singleton(func() *adw.Application {
		return app
	})

	injector.Singleton(func(app *adw.Application) *imgutil.ImgUtil {
		return imgutil.NewImgUtil(app.GetApplicationId())
	})

	if code := app.Run(int32(len(os.Args)), os.Args); code > 0 {
		app.Quit()
		os.Exit(int(code))
	}
}

func onActivate(_ gio.Application) {
	window := ui.NewWindow(app)
	window.Present()
	window.MaybePresentWelcome()
	window.ConnectCloseRequest(new(func(w gtk.Window) bool {
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

func onCommandLine(app gio.Application, ptr uintptr) int32 {
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
		router.Navigate(args[1])
	}
	return 0
}
