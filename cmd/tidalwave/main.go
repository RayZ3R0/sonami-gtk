package main

import (
	"log/slog"
	"os"
	"strings"

	"codeberg.org/dergs/tidalwave/internal/g"
	_ "codeberg.org/dergs/tidalwave/internal/icons"
	"codeberg.org/dergs/tidalwave/internal/player"
	"codeberg.org/dergs/tidalwave/internal/router"
	"codeberg.org/dergs/tidalwave/internal/secrets"
	_ "codeberg.org/dergs/tidalwave/internal/styles"
	"codeberg.org/dergs/tidalwave/internal/ui"
	"codeberg.org/dergs/tidalwave/pkg/schwifty/tracking"
	"codeberg.org/dergs/tidalwave/pkg/tidalapi"
	"codeberg.org/dergs/tidalwave/pkg/utils/imgutil"
	"github.com/infinytum/injector"
	"github.com/jwijenbergh/puregotk/v4/adw"
	"github.com/jwijenbergh/puregotk/v4/gio"
)

func init() {
	slog.SetLogLoggerLevel(slog.LevelDebug)
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

	if code := app.Run(len(os.Args), os.Args); code > 0 {
		app.Quit()
		os.Exit(code)
	}
}

func onActivate(_ gio.Application) {
	window := ui.NewWindow(app)
	window.Present()
	go tracking.LogAliveWidgets()
}

var isActive bool

func onCommandLine(app gio.Application, ptr uintptr) int {
	cmd := gio.ApplicationCommandLineNewFromInternalPtr(ptr)
	args := cmd.GetArguments(nil)
	if !isActive {
		app.Activate()
		isActive = true
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
