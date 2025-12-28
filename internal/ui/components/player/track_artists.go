package player

import (
	"codeberg.org/dergs/tidalwave/internal/player"
	"codeberg.org/dergs/tidalwave/internal/signals"
	"codeberg.org/dergs/tidalwave/pkg/schwifty"
	"codeberg.org/dergs/tidalwave/pkg/schwifty/state"
	. "codeberg.org/dergs/tidalwave/pkg/schwifty/syntax"
)

var artistsState = state.New("")

func init() {
	player.OnTrackChanged.On(func(trackInfo player.TrackInformation) bool {
		artistsState.SetValue(trackInfo.ArtistNames())
		return signals.Continue
	})
}

func trackArtists() schwifty.Label {
	return Label("").
		BindText(artistsState).
		FontSize(16).
		FontWeight(700).
		LineHeight(1.2).
		Color("#1C71D8").
		TextDecoration("underline")
}
