package media_card

import (
	"codeberg.org/puregotk/puregotk/v4/gtk"
	"github.com/RayZ3R0/sonami-gtk/internal/router"
	"github.com/RayZ3R0/sonami-gtk/pkg/schwifty"
	. "github.com/RayZ3R0/sonami-gtk/pkg/schwifty/syntax"
	v2 "github.com/RayZ3R0/sonami-gtk/pkg/tidalapi/models/v2"
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
	var coverUrl string
	if deeplink.Id == "tidal://my-collection/tracks" {
		coverUrl = "https://tidal.com/assets/my-tracks-DTG3pLQW.png"
	}

	return newDeeplink(deeplink.URL, deeplink.Title, coverUrl)
}
