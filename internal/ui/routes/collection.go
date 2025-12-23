package routes

import (
	"codeberg.org/dergs/tidalwave/internal/router"
	"github.com/diamondburned/gotk4/pkg/gtk/v4"
)

func init() {
	router.Register("collection", func(params router.Params) *router.Response {
		label := gtk.NewLabel("Collection goes here")
		return &router.Response{
			PageTitle: "My Collection",
			View:      label,
		}
	})
}
