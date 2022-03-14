package home

import (
	"os"

	"github.com/Pauloo27/mail-notifier/gui/utils"
	"github.com/Pauloo27/mail-notifier/socket/client"
	"github.com/gotk3/gotk3/gtk"
)

func createHeader(c *client.Client) *gtk.HeaderBar {
	container, err := gtk.HeaderBarNew()
	utils.HandleError(err)

	sc, err := container.GetStyleContext()
	utils.HandleError(err)

	sc.AddClass("titlebar")

	titleLbl, err := gtk.LabelNew("Your inboxes:")
	utils.HandleError(err)

	settingsBtn, err := gtk.ButtonNewFromIconName("open-menu", gtk.ICON_SIZE_BUTTON)
	utils.HandleError(err)

	refreshBtn, err := gtk.ButtonNewFromIconName("view-refresh", gtk.ICON_SIZE_BUTTON)
	utils.HandleError(err)

	closeBtn, err := gtk.ButtonNewFromIconName("window-close", gtk.ICON_SIZE_BUTTON)
	utils.HandleError(err)

	closeBtn.Connect("clicked", func() {
		os.Exit(0)
	})

	refreshBtn.Connect("clicked", func() {
		go func() {
			_ = c.ClearAllInboxCache()
			for i := range c.LastInboxList {
				_, _ = c.FetchUnreadMessagesIn(i)
			}
		}()
	})

	container.PackStart(titleLbl)
	container.PackEnd(closeBtn)
	container.PackEnd(settingsBtn)
	container.PackEnd(refreshBtn)

	return container
}
