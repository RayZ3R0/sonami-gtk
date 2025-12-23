package routes

import (
	"codeberg.org/dergs/tidalwave/internal/router"
	"codeberg.org/dergs/tidalwave/internal/ui/signals"
	"github.com/diamondburned/gotk4-adwaita/pkg/adw"
	"github.com/diamondburned/gotk4/pkg/gtk/v4"
)

func init() {
	router.Register("search", func(params router.Params) *router.Response {
		searchBar := gtk.NewSearchEntry()
		searchBar.SetHExpand(true)
		searchBar.SetMarginEnd(40)
		searchBar.SetPlaceholderText("E.g. Fox Stevenson...")
		router.NavigationComplete.On(func(response *router.Response) bool {
			searchBar.GrabFocus()
			return signals.Continue
		})

		searchPrompt := adw.NewStatusPage()
		searchPrompt.SetIconName("loupe-symbolic")
		searchPrompt.SetTitle("Search")
		searchPrompt.SetDescription("Start typing in the search bar to search for songs, artists, albums or playlists.")

		return &router.Response{
			PageTitle: "Search",
			Toolbar:   searchBar,
			View:      searchPrompt,
		}
	})
}
