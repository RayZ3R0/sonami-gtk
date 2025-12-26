package routes

import (
	"context"
	"log/slog"

	"codeberg.org/dergs/tidalwave/internal/router"
	"codeberg.org/dergs/tidalwave/internal/ui/routes/search"
	"codeberg.org/dergs/tidalwave/internal/ui/signals"
	"codeberg.org/dergs/tidalwave/pkg/tidalapi"
	"github.com/diamondburned/gotk4/pkg/gtk/v4"
	"github.com/go-gst/go-glib/glib"
	"github.com/infinytum/injector"
)

func init() {

	router.Register("search", func(params router.Params) *router.Response {
		scrolledWindow := gtk.NewScrolledWindow()
		scrolledWindow.SetPolicy(gtk.PolicyNever, gtk.PolicyAutomatic)

		searchBar := gtk.NewSearchEntry()
		searchBar.SetHExpand(true)
		searchBar.SetMarginEnd(40)
		searchBar.SetPlaceholderText("E.g. Fox Stevenson...")
		searchBar.ConnectSearchChanged(OnSearch(searchBar, scrolledWindow))
		searchBar.SetSearchDelay(1000)
		router.NavigationComplete.On(func(response *router.Response) bool {
			searchBar.GrabFocus()
			return signals.Continue
		})
		scrolledWindow.SetChild(search.PromptView())

		return &router.Response{
			PageTitle: "Search",
			Toolbar:   searchBar,
			View:      scrolledWindow,
		}
	})
}

var searchIncludes = []string{
	"topHits", "topHits.profileArt", "topHits.coverArt", "topHits.albums.coverArt",
	"topHits.albums.artists",
}

func OnSearch(searchBar *gtk.SearchEntry, scrolledWindow *gtk.ScrolledWindow) func() {
	tidal := injector.MustInject[*tidalapi.TidalAPI]()
	return func() {
		query := searchBar.Text()
		if query == "" {
			scrolledWindow.SetChild(search.PromptView())
			return
		}

		slog.Info("searching", "query", query)
		scrolledWindow.SetChild(search.LoadingView())
		go func() {
			searchResults, err := tidal.OpenAPI.V2.SearchResults.Search(context.Background(), query, searchIncludes...)
			glib.IdleAdd(func() {
				if err != nil {
					slog.Error("search failed", "error", err)
					scrolledWindow.SetChild(search.PromptView())
					return
				}

				scrolledWindow.SetChild(search.TopHits(searchResults))
			})
		}()
	}
}
