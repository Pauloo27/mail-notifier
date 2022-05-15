package message

import (
	"github.com/Pauloo27/mail-notifier/gui/utils"
	"github.com/Pauloo27/mail-notifier/socket/common/types"
	"github.com/gotk3/gotk3/gtk"
)

func Show(msg *types.CachedMailMessage, contentType string) {
	gtk.Init(nil)

	win, err := gtk.WindowNew(gtk.WINDOW_TOPLEVEL)
	utils.HandleError(err)

	win.Connect("destroy", func() {
		gtk.MainQuit()
	})

	mainContainer, err := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 1)
	utils.HandleError(err)

	mainContainer.PackStart(createHeader(msg, win), false, false, 5)
	mainContainer.PackStart(createMessageHeaderPreview(msg, contentType), false, true, 5)
	mainContainer.PackStart(createMessageContentPreview(msg, contentType), false, true, 5)

	win.Add(mainContainer)

	win.SetTitle("Mail notifier - Message")
	win.SetDefaultSize(450, 500)
	win.ShowAll()
	gtk.Main()
}
