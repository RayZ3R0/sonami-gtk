package tidalapi

import (
	"fmt"
	"strings"
)

func ImageURL(id string) string {
	return "https://resources.tidal.com/images/" + strings.ReplaceAll(id, "-", "/") + "/320x320.jpg"
}

func FormatDuration(sec int) string {
	if sec < 0 {
		sec = 0
	}

	const (
		day  = 24 * 3600
		hour = 3600
		min  = 60
	)

	d := sec / day
	sec %= day

	h := sec / hour
	sec %= hour

	m := sec / min
	s := sec % min

	switch {
	case d > 0:
		return fmt.Sprintf("%dd %d:%02d:%02d", d, h, m, s)
	case h > 0:
		return fmt.Sprintf("%d:%02d:%02d", h, m, s)
	default:
		return fmt.Sprintf("%d:%02d", m, s)
	}
}
