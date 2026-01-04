package ui

import (
	"fmt"
	"regexp"

	"github.com/jwijenbergh/puregotk/v4/adw"
	"github.com/jwijenbergh/puregotk/v4/gtk"
)

var (
	Commit, Version string
)

func getVersionPrefix(version string) (prefix string) {
	if version[0] == 'v' {
		if version[1] == '0' {
			prefix = "β "
		}
	} else {
		if version[0] == '0' {
			prefix = "β "
		} else {
			prefix = "v"
		}
	}

	return
}

func getVersionNumber() (version string) {
	version = Version

	if version == "" {
		// If no version is available.
		// This happens if the go binary is run directly, or through the Zed task

		if Commit == "" {
			// If no commit is available.
			// This happens if the app is run with `go run ./cmd/tidalwave`.

			version = "local"
		} else if ok, _ := regexp.MatchString(`^.*-\d+-g[0-9a-f]{7}$`, Commit); ok {
			// If the commit matches the git describe scheme.
			// This happens when the code is run locally, through a project-provided tool, and a tag is available.
			// We can therefore display a git version, in the `v1.0.0 (r1.0123abc)` scheme.

			regex := regexp.MustCompile(`^(.*)-(\d+)-g([0-9a-f]{7})`)
			parts := regex.FindStringSubmatch(Commit)

			_, tag, offset, commitHash := parts[0], parts[1], parts[2], parts[3]
			prefix := getVersionPrefix(tag)

			version = fmt.Sprintf("%s%s (r%s.%s)", prefix, tag, offset, commitHash)
		} else {
			// If the commit is not a git describe scheme.
			// This happens when the code is run locally, through a project-provided tool, and no tag is available.
			// We can therefore display a commit hash, with the alpha symbol.
			version = fmt.Sprintf("α git+%s", Commit[:7])
		}
	} else {
		// If a version has been provided.
		// We can therefore display the version, and the commit hash for clarity.

		commit := Commit
		if commit == "" {
			commit = "dirty!"
		} else {
			commit = Commit[:7]
		}

		prefix := getVersionPrefix(Version)

		version = fmt.Sprintf("%s%s (%s)", prefix, Version, commit)
	}

	return
}

func (w *Window) PresentAbout() {

	about := adw.NewAboutDialog()
	about.SetApplicationIcon("logo")
	about.SetApplicationName("Tidal Wave")
	about.SetVersion(getVersionNumber())
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
