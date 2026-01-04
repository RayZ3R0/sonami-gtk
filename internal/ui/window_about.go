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
		"Nila The Dragon https://github.com/NilaTheDragon",
		"Dråfølin https://github.com/Drafolin",
	})
	about.SetCopyright("© 2026 Nila The Dragon")
	about.SetWebsite("https://codeberg.org/dergs/TidalWave")
	about.SetIssueUrl("https://codeberg.org/dergs/TidalWave/issues")

	about.AddLegalSection("GStreamer Bindings (go-gst/go-gst)", "© 2020 https://github.com/go-gst/go-gst", gtk.LicenseLgpl30Value, "")
	about.AddLegalSection("DBus Client (godbus/dbus)", "© 2020 Georg Reinke https://github.com/godbus/dbus", gtk.LicenseBsdValue, "")
	about.AddLegalSection("UUID (google/uuid)", "© 2009, 2014 Google Inc. https://github.com/google/uuid", gtk.LicenseBsd3Value, "")
	about.AddLegalSection("Dependency Injector (infinytum/injector)", "© 2022 Infinytum https://github.com/infinytum/injector", gtk.LicenseUnknownValue, "")
	about.AddLegalSection("System Locale Detection (jeandeaual/go-locale)", "© 2020 Alexis Jeandeau https://github.com/jeandeaual/go-locale", gtk.LicenseMitX11Value, "")
	about.AddLegalSection("GTK4 / Libadwaita Bindings (jwijenbergh/puregotk)", "© 2022 Kyle McGough https://github.com/jwijenbergh/puregotk", gtk.LicenseMitX11Value, "")
	about.AddLegalSection("ISO8601 Duration Parser (osodev/duration)", "© 2023 Jeroen Wijenbergh https://github.com/sosodev/duration", gtk.LicenseMitX11Value, "")
	about.AddLegalSection("QR Code Generator (yeqown/go-qrcode)", "© 2018 yeqown https://github.com/yeqown/go-qrcode", gtk.LicenseMitX11Value, "")
	about.AddLegalSection("Keyring (zalando/go-keyring)", "© 2016 Zalando SE https://github.com/zalando/go-keyring", gtk.LicenseMitX11Value, "")

	about.Present(&w.Widget)
	about.Unref()
}
