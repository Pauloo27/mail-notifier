package home

import (
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

	refreshBtn.Connect("clicked", func() {
		go func() {
			_ = c.ClearAllInboxCache()
			for i := range c.LastInboxList {
				_, _ = c.FetchUnreadMessagesIn(i)
			}
		}()
	})

	container.PackStart(titleLbl)
	container.PackEnd(settingsBtn)
	container.PackEnd(refreshBtn)

	return container
}
