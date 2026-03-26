package media_card

import (
	"codeberg.org/puregotk/puregotk/v4/glib"
	"github.com/RayZ3R0/sonami-gtk/pkg/schwifty"
	v2 "github.com/RayZ3R0/sonami-gtk/pkg/tidalapi/models/v2"
)

func NewMixGeneric(id string, title string, subtitle string, coverUrl string) schwifty.Button {
	return Card(
		title,
		SubTitle(subtitle).Lines(2),
		coverUrl,
	).ActionName("win.route.playlist").ActionTargetValue(glib.NewVariantString(id))
}

func NewLegacyMix(mix *v2.MixItemData) schwifty.Button {
	return NewMixGeneric(mix.Id, mix.TitleTextInfo.Text, mix.SubtitleTextInfo.Text, mix.MixImages[0].URL)
}
