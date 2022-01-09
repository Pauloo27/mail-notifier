package home

import (
	"github.com/Pauloo27/mail-notifier/gui/utils"
	"github.com/Pauloo27/mail-notifier/socket/client"
	"github.com/gotk3/gotk3/gtk"
)

func Show(c *client.Client) {
	gtk.Init(nil)

	win, err := gtk.WindowNew(gtk.WINDOW_TOPLEVEL)
	utils.HandleError(err)

	win.SetPosition(gtk.WIN_POS_MOUSE)

	win.Connect("destroy", func() {
		gtk.MainQuit()
	})

	mainContainer, err := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 1)
	utils.HandleError(err)

	mainContainer.PackStart(createHeader(c), false, false, 5)
	mainContainer.PackStart(createInboxList(c), true, true, 5)

	win.Add(mainContainer)

	win.SetTitle("Mail notifier")
	win.SetDefaultSize(450, 500)
	win.ShowAll()
	gtk.Main()
}
