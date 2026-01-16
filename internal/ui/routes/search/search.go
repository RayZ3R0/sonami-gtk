package search

import (
	"codeberg.org/dergs/tonearm/pkg/schwifty/bindings/adw"
	. "codeberg.org/dergs/tonearm/pkg/schwifty/syntax"
)

func PromptView() adw.StatusPage {
	return StatusPage().
		IconName("loupe-symbolic").
		Title("Search").
		Description("Start typing in the search bar to search for songs, artists, albums or playlists.")
}

func LoadingView() adw.Clamp {
	return Clamp().
		MaximumSize(50).
		Child(Spinner())
}
