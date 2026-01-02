package ui

import (
	"github.com/jwijenbergh/puregotk/v4/adw"
	"github.com/jwijenbergh/puregotk/v4/gtk"
)

func (w *Window) PresentAbout() {
	about := adw.NewAboutDialog()
	about.SetApplicationIcon("logo")
	about.SetApplicationName("Tidal Wave")
	about.SetVersion("git")
	about.SetLicenseType(gtk.LicenseGpl30Value)
	about.SetDevelopers([]string{
		"Nila The Dragon",
	})

	about.Present(&w.Widget)
}
