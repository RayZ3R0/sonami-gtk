package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"

	_ "codeberg.org/dergs/tidalwave/internal/icons"
	"codeberg.org/dergs/tidalwave/internal/ui"
	"codeberg.org/dergs/tidalwave/pkg/tidalapi"
	"github.com/diamondburned/gotk4-adwaita/pkg/adw"
	"github.com/diamondburned/gotk4/pkg/core/glib"
	"github.com/diamondburned/gotkit/app"
	"github.com/diamondburned/gotkit/components/prefui"
	"github.com/infinytum/injector"
)

func init() {
	adw.Init()
	slog.SetLogLoggerLevel(slog.LevelDebug)
}

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	app := app.New(context.Background(), "org.codeberg.dergs.tidalwave", "TIDAL for GNOME")
	injector.Singleton(func() context.Context {
		return app.Context()
	})

	var window *ui.Window
	app.AddJSONActions(map[string]interface{}{
		"app.about":       func() { window.PresentAbout() },
		"app.preferences": func() { prefui.ShowDialog(app.Context()) },
		"app.quit":        func() { app.Quit() },
	})
	app.AddActionShortcuts(map[string]string{
		"<Ctrl>Q": "app.quit",
	})
	app.ConnectActivate(func() {
		window = ui.NewWindow(app.Context())
		window.Present()
	})

	go func() {
		<-ctx.Done()
		glib.IdleAdd(app.Quit)
	}()
	injector.Singleton(func() *tidalapi.TidalAPI {
		countryCode, err := tidalapi.FetchCountryCode()
		if err != nil {
			slog.Error("Failed to fetch country code, defaulting to WW", err)
			countryCode = "WW"
		}
		slog.Info("Discovered country code", "countryCode", countryCode)
		return tidalapi.NewClient(countryCode)
	})

	if code := app.Run(os.Args); code > 0 {
		cancel()
		os.Exit(code)
	}
}
