package media_card

import (
	"codeberg.org/dergs/tonearm/internal/router"
	"codeberg.org/dergs/tonearm/pkg/schwifty"
	. "codeberg.org/dergs/tonearm/pkg/schwifty/syntax"
	v2 "codeberg.org/dergs/tonearm/pkg/tidalapi/models/v2"
	"github.com/jwijenbergh/puregotk/v4/gtk"
)

func newDeeplink(url string, name string, coverUrl string) schwifty.Button {
	return Card(
		name,
		HStack(),
		coverUrl,
	).ConnectClicked(func(b gtk.Button) {
		router.Navigate(url)
	})
}

func NewLegacyDeeplink(deeplink *v2.DeepLinkItemData) schwifty.Button {
	return newDeeplink(deeplink.URL, deeplink.Title, "")
}
