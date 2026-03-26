package routes

import (
	"context"
	"encoding/json"
	"log/slog"

	"codeberg.org/puregotk/puregotk/v4/gtk"
	"github.com/RayZ3R0/sonami-gtk/internal/gettext"
	"github.com/RayZ3R0/sonami-gtk/internal/localdb"
	"github.com/RayZ3R0/sonami-gtk/internal/router"
	"github.com/RayZ3R0/sonami-gtk/internal/ui/components"
	. "github.com/RayZ3R0/sonami-gtk/pkg/schwifty/syntax"
	"github.com/RayZ3R0/sonami-gtk/pkg/tidalapi"
	v2 "github.com/RayZ3R0/sonami-gtk/pkg/tidalapi/models/v2"
	"github.com/infinytum/injector"
)

const homeFeedCacheKey = "home_feed_static"

var homeLogger = slog.With("module", "routes/home")

func init() {
	router.Register("home", Home)
}

func Home() *router.Response {
	tidal := injector.MustInject[*tidalapi.TidalAPI]()

	// Try to load from cache first
	var homeFeed *v2.Page
	if cached, err := localdb.GetPageFeed(homeFeedCacheKey); err == nil && cached != nil {
		var page v2.Page
		if err := json.Unmarshal(cached.Data, &page); err == nil {
			homeFeed = &page
			homeLogger.Debug("loaded home feed from cache")

			// Refresh cache in background (non-blocking)
			go refreshHomeFeedCache(tidal)
		} else {
			homeLogger.Warn("failed to deserialize cached home feed, fetching fresh", "error", err)
			// Delete corrupted cache
			_ = localdb.DeletePageFeed(homeFeedCacheKey)
		}
	}

	// If no cache, fetch from API (blocking)
	if homeFeed == nil {
		homeLogger.Debug("fetching home feed from API")
		feed, err := tidal.V2.Home.Feed.Static(context.Background())
		if err != nil {
			return router.FromError(gettext.Get("Home"), err)
		}
		homeFeed = feed

		// Cache the response
		go cacheHomeFeed(feed)
	}

	body := VStack().Spacing(25).VMargin(20)
	for _, item := range homeFeed.Items {
		body = body.Append(components.ForPageItem(item))
	}

	return &router.Response{
		PageTitle: gettext.Get("Home"),
		View: ScrolledWindow().
			Child(
				components.MainContent(
					body,
				),
			).
			Policy(gtk.PolicyNeverValue, gtk.PolicyAutomaticValue),
	}
}

// refreshHomeFeedCache fetches fresh data and updates the cache
func refreshHomeFeedCache(tidal *tidalapi.TidalAPI) {
	feed, err := tidal.V2.Home.Feed.Static(context.Background())
	if err != nil {
		homeLogger.Debug("background refresh failed", "error", err)
		return
	}
	cacheHomeFeed(feed)
	homeLogger.Debug("home feed cache refreshed in background")
}

// cacheHomeFeed serializes and stores the home feed
func cacheHomeFeed(feed *v2.Page) {
	data, err := json.Marshal(feed)
	if err != nil {
		homeLogger.Warn("failed to serialize home feed for cache", "error", err)
		return
	}
	if err := localdb.SetPageFeed(homeFeedCacheKey, data); err != nil {
		homeLogger.Warn("failed to cache home feed", "error", err)
	}
}
