package main

import (
	"log/slog"
	"os"

	"codeberg.org/dergs/tidalwave/internal/g"
	_ "codeberg.org/dergs/tidalwave/internal/icons"
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
	app = adw.NewApplication("org.codeberg.dergs.tidalwave", gio.GApplicationFlagsNoneValue)
	defer app.Unref()
	app.ConnectActivate(g.Ptr(onActivate))
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
