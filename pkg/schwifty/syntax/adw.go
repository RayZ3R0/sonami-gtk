package syntax

import (
	"codeberg.org/dergs/tonearm/pkg/schwifty/bindings"
	adwbindings "codeberg.org/dergs/tonearm/pkg/schwifty/bindings/adw"
	gtkbindings "codeberg.org/dergs/tonearm/pkg/schwifty/bindings/gtk"
	"codeberg.org/dergs/tonearm/pkg/schwifty/callback"
	"github.com/jwijenbergh/puregotk/v4/adw"
	"github.com/jwijenbergh/puregotk/v4/gtk"
)

func AlertDialog(heading string, body string) adwbindings.AlertDialog {
	return managedWidget("AlertDialog", func() *adw.AlertDialog {
		dialog := adw.NewAlertDialog(heading, body)
		dialog.ConnectCloseAttempt(&callback.AlertDialogCloseAttempt)
		return dialog
	})
}

func Clamp() adwbindings.Clamp {
	return managedWidget("Clamp", func() *adw.Clamp {
		return adw.NewClamp()
	})
}

func EntryRow() adwbindings.EntryRow {
	return managedWidget("EntryRow", func() *adw.EntryRow {
		return adw.NewEntryRow()
	})
}

func HeaderBar() adwbindings.HeaderBar {
	return managedWidget("HeaderBar", func() *adw.HeaderBar {
		return adw.NewHeaderBar()
	})
}

func PasswordEntryRow() adwbindings.PasswordEntryRow {
	return managedWidget("PasswordEntryRow", func() *adw.PasswordEntryRow {
		return adw.NewPasswordEntryRow()
	})
}

func PreferencesDialog(pages ...any) adwbindings.PreferencesDialog {
	return managedWidget("PreferencesDialog", func() *adw.PreferencesDialog {
		dialog := adw.NewPreferencesDialog()
		for _, page := range pages {
			dialog.Add(bindings.ResolveTo[*adw.PreferencesPage, adwbindings.PreferencesPage](page))
		}
		return dialog
	})
}

func PreferencesGroup(children ...any) adwbindings.PreferencesGroup {
	return managedWidget("PreferencesGroup", func() *adw.PreferencesGroup {
		group := adw.NewPreferencesGroup()
		for _, child := range children {
			group.Add(gtkbindings.ResolveWidget(child))
		}
		return group
	})
}

func PreferencesPage(groups ...any) adwbindings.PreferencesPage {
	return managedWidget("PreferencesPage", func() *adw.PreferencesPage {
		page := adw.NewPreferencesPage()
		for _, group := range groups {
			page.Add(bindings.ResolveTo[*adw.PreferencesGroup, adwbindings.PreferencesGroup](group))
		}
		return page
	})
}

func SpinRow(adjustment *gtk.Adjustment, climbRate float64, digits uint) adwbindings.SpinRow {
	return managedWidget("SpinRow", func() *adw.SpinRow {
		return adw.NewSpinRow(adjustment, climbRate, digits)
	})
}

func StatusPage() adwbindings.StatusPage {
	return managedWidget("StatusPage", func() *adw.StatusPage {
		return adw.NewStatusPage()
	})
}

func SwitchRow() adwbindings.SwitchRow {
	return managedWidget("SwitchRow", func() *adw.SwitchRow {
		return adw.NewSwitchRow()
	})
}

func ViewStack(children ...any) adwbindings.ViewStack {
	return managedWidget("ViewStack", func() *adw.ViewStack {
		viewStack := adw.NewViewStack()
		for _, child := range children {
			viewStack.Add(gtkbindings.ResolveWidget(child))
		}
		return viewStack
	})
}

func WindowTitle(title string, subtitle string) adwbindings.WindowTitle {
	return managedWidget("WindowTitle", func() *adw.WindowTitle {
		return adw.NewWindowTitle(title, subtitle)
	})
}

func WrapBox(children ...any) adwbindings.WrapBox {
	return managedWidget("WrapBox", func() *adw.WrapBox {
		box := adw.NewWrapBox()
		for _, child := range children {
			box.Append(gtkbindings.ResolveWidget(child))
		}
		return box
	})
}
