package v2

import (
	"regexp"

	"github.com/RayZ3R0/sonami-gtk/pkg/sonami"
	v2 "github.com/RayZ3R0/sonami-gtk/pkg/tidalapi/models/v2"
)

var wimpLinkRegex = regexp.MustCompile(`\[wimpLink\s+artistId="[^"]*"\]([^[]*)\[/wimpLink\]`)

type Artist struct {
	ArtistInfo
	ArtistPage v2.ArtistPage
}

func (a *Artist) Description() string {
	return wimpLinkRegex.ReplaceAllString(a.ArtistPage.Header.Biography.Text, "$1")
}

func (a Artist) FollowerCount() int {
	return a.ArtistPage.Header.FollowersAmount
}

func NewArtist(artist v2.ArtistPage) sonami.Artist {
	return &Artist{ArtistInfo{*artist.Item.Data.Artist}, artist}
}
