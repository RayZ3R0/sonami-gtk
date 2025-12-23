package tidalapi

import "strings"

func ImageURL(id string) string {
	return "https://resources.tidal.com/images/" + strings.ReplaceAll(id, "-", "/") + "/320x320.jpg"
}
