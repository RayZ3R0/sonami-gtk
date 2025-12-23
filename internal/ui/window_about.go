package ui

import (
	"github.com/diamondburned/gotk4-adwaita/pkg/adw"
	"github.com/diamondburned/gotk4/pkg/gtk/v4"
)

func (w *Window) PresentAbout() {
	about := adw.NewAboutDialog()
	about.SetApplicationIcon("logo")
	about.SetApplicationName("Tidal Wave")
	about.SetVersion("git")
	about.SetLicenseType(gtk.LicenseGPL30)
	about.SetDevelopers([]string{
		"Nila The Dragon",
	})

	about.Present(w)
}
