package navigation

import (
	"codeberg.org/dergs/tonearm/internal/signals"
)

type Path string

var (
	PathPlayer Path = "player"
	PathLyrics Path = "lyrics"
	PathQueue  Path = "queue"
)

// SidebarNavigation is a signal that is emitted when the router starts or completes a navigation in the sidebar.
// If Completed is true, the Result parameter is always set.
var Navigation = signals.NewStatelessSignal[Path]()
