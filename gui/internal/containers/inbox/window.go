package inbox

import (
	"github.com/Pauloo27/mail-notifier/gui/utils"
	"github.com/Pauloo27/mail-notifier/core/provider"
	"github.com/gotk3/gotk3/gtk"
)

func Show(mail provider.MailProvider, messages []provider.MailMessage) {
	gtk.Init(nil)

	win, err := gtk.WindowNew(gtk.WINDOW_TOPLEVEL)
	utils.HandleError(err)

	win.Connect("destroy", func() {
		gtk.MainQuit()
	})

	mainContainer, err := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 1)
	utils.HandleError(err)

	mainContainer.PackStart(createHeader(mail, win), false, false, 5)
	mainContainer.PackStart(createMessageList(mail, messages), true, true, 5)

	win.Add(mainContainer)

	win.SetTitle("Mail notifier - Inbox")
	win.SetDefaultSize(450, 500)
	win.ShowAll()
	gtk.Main()
}
