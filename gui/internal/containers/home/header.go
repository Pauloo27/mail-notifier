package home

import (
	"github.com/Pauloo27/mail-notifier/gui/utils"
	"github.com/gotk3/gotk3/gtk"
)

func createHeader() *gtk.HeaderBar {
	container, err := gtk.HeaderBarNew()
	utils.HandleError(err)

	titleLbl, err := gtk.LabelNew("Your inboxes:")
	utils.HandleError(err)

	settingsBtn, err := gtk.ButtonNewFromIconName("open-menu", gtk.ICON_SIZE_BUTTON)
	utils.HandleError(err)

	container.PackStart(titleLbl)
	container.PackEnd(settingsBtn)

	return container
}
