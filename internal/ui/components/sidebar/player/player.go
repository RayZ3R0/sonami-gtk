package player

import (
	"log/slog"

	"codeberg.org/puregotk/puregotk/v4/gtk"
	"github.com/RayZ3R0/sonami-gtk/internal/player"
	"github.com/RayZ3R0/sonami-gtk/internal/signals"
	"github.com/RayZ3R0/sonami-gtk/pkg/schwifty"
	"github.com/RayZ3R0/sonami-gtk/pkg/schwifty/state"
	. "github.com/RayZ3R0/sonami-gtk/pkg/schwifty/syntax"
)

var logger = slog.With("module", "ui/components", "component", "player")

var (
	isControllableState = state.NewStateful(false)
)

func init() {
	player.PlaybackStateChanged.On(func(ps *player.PlaybackState) bool {
		isControllableState.SetValue(!ps.Loading)
		return signals.Continue
	})
}

func NewPlayer() schwifty.Box {
	return VStack(
		PlayingFrom(),
		trackCover(),
		trackInfo(),
		actionRow(),
		trackTimeline(),
		controls(),
	).
		WithCSSClass("player").
		HPadding(24).
		Spacing(25).
		VAlign(gtk.AlignCenterValue)
}
