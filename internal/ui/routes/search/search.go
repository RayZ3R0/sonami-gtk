package search

import (
	"github.com/RayZ3R0/sonami-gtk/internal/gettext"
	"github.com/RayZ3R0/sonami-gtk/pkg/schwifty"
	. "github.com/RayZ3R0/sonami-gtk/pkg/schwifty/syntax"
)

func PromptView() schwifty.StatusPage {
	return StatusPage().
		IconName("loupe-symbolic").
		Title(gettext.Get("Search")).
		Description(gettext.Get("Start typing in the search bar to search for songs, artists, albums or playlists"))
}

func LoadingView() schwifty.Clamp {
	return Clamp().
		MaximumSize(50).
		Child(Spinner())
}
