package routes

import (
	"context"
	"log/slog"
	"time"

	"codeberg.org/puregotk/puregotk/v4/gtk"
	"github.com/RayZ3R0/sonami-gtk/internal/gettext"
	"github.com/RayZ3R0/sonami-gtk/internal/router"
	"github.com/RayZ3R0/sonami-gtk/internal/ui/components"
	"github.com/RayZ3R0/sonami-gtk/internal/ui/routes/search"
	"github.com/RayZ3R0/sonami-gtk/pkg/schwifty/state"
	. "github.com/RayZ3R0/sonami-gtk/pkg/schwifty/syntax"
	"github.com/RayZ3R0/sonami-gtk/pkg/tidalapi"
	"github.com/infinytum/injector"
)

func init() {
	router.Register("search", func() *router.Response {
		scrollChildState := state.NewStateful[any](search.PromptView())
		searchState := state.NewStateful(false)
		searchHandler := OnSearch(scrollChildState)

		return &router.Response{
			PageTitle: gettext.Get("Search"),
			Toolbar: Clamp().
				Orientation(gtk.OrientationHorizontalValue).
				MaximumSize(500).
				Child(
					SearchEntry().
						HExpand(true).
						PlaceholderText(gettext.Get("For example, Fox Stevenson")).
						SearchDelay(1000).
						ConnectActivate(func(se gtk.SearchEntry) {
							searchState.SetValue(true)
							time.AfterFunc(time.Second, func() {
								searchState.SetValue(false)
							})
							searchHandler(se)
						}).
						ConnectMap(func(w gtk.Widget) {
							w.GrabFocus()
						}).
						ConnectSearchChanged(func(se gtk.SearchEntry) {
							if searchState.Value() && se.GetText() != "" {
								return
							}
							searchHandler(se)
						}),
				),
			View: ScrolledWindow().
				BindChild(scrollChildState).
				Policy(gtk.PolicyNeverValue, gtk.PolicyAutomaticValue),
		}
	})
}

var searchIncludes = []string{
	"topHits", "topHits.profileArt", "topHits.coverArt", "topHits.albums.coverArt",
	"topHits.albums.artists", "topHits.artists",
}

func OnSearch(scrollChildState *state.State[any]) func(gtk.SearchEntry) {
	tidal := injector.MustInject[*tidalapi.TidalAPI]()
	return func(searchBar gtk.SearchEntry) {
		query := searchBar.GetText()
		if query == "" {
			scrollChildState.SetValue(search.PromptView())
			return
		}
		scrollChildState.SetValue(search.LoadingView())
		go func() {
			searchResults, err := tidal.OpenAPI.V2.SearchResults.Search(context.Background(), query, searchIncludes...)
			if err != nil {
				slog.Error("search failed", "error", err)
				scrollChildState.SetValue(search.PromptView())
				return
			}
			scrollChildState.SetValue(
				components.MainContent(
					search.TopHits(searchResults),
				),
			)
		}()
	}
}
