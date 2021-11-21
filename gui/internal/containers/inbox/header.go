package inbox

import (
	"fmt"

	"github.com/Pauloo27/mail-notifier/gui/utils"
	"github.com/Pauloo27/mail-notifier/internal/provider"
	"github.com/gotk3/gotk3/gtk"
)

func createHeader(mail provider.MailProvider) *gtk.HeaderBar {
	container, err := gtk.HeaderBarNew()
	utils.HandleError(err)

	titleLbl, err := gtk.LabelNew(fmt.Sprintf("%s messages:", mail.GetAddress()))
	utils.HandleError(err)

	container.PackStart(titleLbl)

	return container
}
