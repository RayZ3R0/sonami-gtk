package media_card

import (
	"codeberg.org/dergs/tidalwave/internal/router"
	"codeberg.org/dergs/tidalwave/pkg/schwifty"
	. "codeberg.org/dergs/tidalwave/pkg/schwifty/syntax"
	v2 "codeberg.org/dergs/tidalwave/pkg/tidalapi/models/v2"
	"github.com/jwijenbergh/puregotk/v4/gtk"
)

func newDeeplink(url string, name string, coverUrl string) schwifty.Button {
	return Card(
		name,
		HStack(),
		coverUrl,
	).ConnectClicked(func(b gtk.Button) {
		router.Navigate(url, nil)
	})
}

func NewLegacyDeeplink(deeplink *v2.DeepLinkItemData) schwifty.Button {
	return newDeeplink(deeplink.URL, deeplink.Title, "")
}
