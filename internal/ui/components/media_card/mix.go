package media_card

import (
	"codeberg.org/dergs/tidalwave/pkg/schwifty"
	v2 "codeberg.org/dergs/tidalwave/pkg/tidalapi/models/v2"
	"github.com/jwijenbergh/puregotk/v4/glib"
)

func newMix(id string, title string, subtitle string, coverUrl string) schwifty.Button {
	return Card(
		title,
		SubTitle(subtitle).Lines(2),
		coverUrl,
	).ActionName("win.route.playlist").ActionTargetValue(glib.NewVariantString(id))
}

func NewLegacyMix(mix *v2.MixItemData) schwifty.Button {
	return newMix(mix.Id, mix.TitleTextInfo.Text, mix.SubtitleTextInfo.Text, mix.MixImages[0].URL)
}
