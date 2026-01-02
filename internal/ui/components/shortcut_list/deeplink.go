package shortcut_list

import (
	"codeberg.org/dergs/tidalwave/pkg/schwifty"
	v2 "codeberg.org/dergs/tidalwave/pkg/tidalapi/models/v2"
)

func newDeepLink(id string, name string, coverUrl string) schwifty.Button {
	return NewShortcut(
		name,
		"",
		coverUrl,
	)
}

func NewLegacyDeepLink(deepLink *v2.DeepLinkItemData) schwifty.Button {
	return newDeepLink(deepLink.Id, deepLink.Title, "")
}
