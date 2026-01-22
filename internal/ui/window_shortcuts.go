package ui

import (
	. "codeberg.org/dergs/tonearm/pkg/schwifty/syntax"
)

func (w *Window) PresentShortcuts() {
	ShortcutsDialog(
		ShortcutsSection(
			ShortcutsItemFromAction("Close", "win.close"),
			ShortcutsItemFromAction("Quit", "app.quit"),
			ShortcutsItemFromAction("Side Pane", "win.toggle-sidebar"),
			ShortcutsItemFromAction("Main Menu", "win.main-menu"),
			ShortcutsItemFromAction("Keyboard Shortcuts", "app.shortcuts"),
			ShortcutsItemFromAction("Preferences", "app.preferences"),
		).Title("Basic Shortcuts"),
		ShortcutsSection(
			ShortcutsItemFromAction("Back", "win.navigate-back"),
			ShortcutsItemFromAction("Search", "win.search"),
		).Title("Navigation"),
	).Present(w)
}
