package inbox

import (
	"fmt"

	"github.com/Pauloo27/mail-notifier/gui/utils"
	"github.com/Pauloo27/mail-notifier/internal/provider"
	"github.com/gotk3/gotk3/gtk"
)

func createHeader(mail provider.MailProvider, win *gtk.Window) *gtk.HeaderBar {
	container, err := gtk.HeaderBarNew()
	utils.HandleError(err)

	titleLbl, err := gtk.LabelNew(fmt.Sprintf("%s messages:", mail.GetAddress()))
	utils.HandleError(err)

	closeBtn, err := gtk.ButtonNewFromIconName("go-previous", gtk.ICON_SIZE_BUTTON)
	utils.HandleError(err)

	closeBtn.Connect("clicked", func() {
		win.Destroy()
	})

	container.PackStart(closeBtn)
	container.PackStart(titleLbl)

	return container
}
