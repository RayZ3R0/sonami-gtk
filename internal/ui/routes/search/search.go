package search

import (
	"codeberg.org/dergs/tidalwave/pkg/schwifty"
	. "codeberg.org/dergs/tidalwave/pkg/schwifty/syntax"
)

func PromptView() schwifty.StatusPage {
	return StatusPage().
		IconName("loupe-symbolic").
		Title("Search").
		Description("Start typing in the search bar to search for songs, artists, albums or playlists.")
}

func LoadingView() schwifty.Clamp {
	return Clamp().
		MaximumSize(50).
		Child(Spinner())
}
