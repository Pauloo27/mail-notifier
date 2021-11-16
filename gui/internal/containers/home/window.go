package home

import (
	"github.com/Pauloo27/mail-notifier/gui/utils"
	"github.com/gotk3/gotk3/gtk"
)

func Show() {
	gtk.Init(nil)

	win, err := gtk.WindowNew(gtk.WINDOW_TOPLEVEL)
	utils.HandleError(err)

	win.Connect("destroy", func() {
		gtk.MainQuit()
	})

	mainContainer, err := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 1)
	utils.HandleError(err)

	mainContainer.PackStart(createHeader(), false, false, 5)
	mainContainer.PackStart(createInboxList(), true, true, 5)

	win.Add(mainContainer)

	win.SetTitle("Mail notifier")
	win.SetDefaultSize(300, 500)
	win.ShowAll()
	gtk.Main()
}
