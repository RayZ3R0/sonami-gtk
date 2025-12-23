package routes

import (
	"codeberg.org/dergs/tidalwave/internal/router"
	"github.com/diamondburned/gotk4/pkg/gtk/v4"
)

func init() {
	router.Register("explore", func(params router.Params) *router.Response {
		label := gtk.NewLabel("Explore goes here")
		return &router.Response{
			PageTitle: "Explore",
			View:      label,
		}
	})
}
