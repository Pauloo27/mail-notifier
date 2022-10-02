package message

import (
	"github.com/Pauloo27/mail-notifier/gui/utils"
	"github.com/Pauloo27/mail-notifier/socket/common/types"
	"github.com/gotk3/gotk3/gtk"
)

func createHeader(msg *types.CachedMailMessage, win *gtk.Window) *gtk.HeaderBar {
	container, err := gtk.HeaderBarNew()
	utils.HandleError(err)

	sc, err := container.GetStyleContext()
	utils.HandleError(err)

	sc.AddClass("titlebar")

	titleLbl, err := gtk.LabelNew(utils.AddEllipsis(msg.Subject, 30))
	utils.HandleError(err)

	titleLbl.SetTooltipText(msg.Subject)

	closeBtn, err := gtk.ButtonNewFromIconName("go-previous", gtk.ICON_SIZE_BUTTON)
	utils.HandleError(err)

	closeBtn.Connect("clicked", func() {
		win.Destroy()
	})

	container.PackStart(closeBtn)
	container.PackStart(titleLbl)

	return container
}
