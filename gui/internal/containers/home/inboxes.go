package home

import (
	"strconv"

	"github.com/Pauloo27/mail-notifier/gui/internal/containers/inbox"
	"github.com/Pauloo27/mail-notifier/gui/utils"
	"github.com/Pauloo27/mail-notifier/socket/client"
	"github.com/Pauloo27/mail-notifier/socket/common/types"
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
)

func createInboxItem(box *types.Inbox, messages *types.CachedUnreadMessages) *gtk.Box {
	container, err := gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 5)
	utils.HandleError(err)

	address := box.Address

	emailLbl, err := gtk.LabelNew(address)
	utils.HandleError(err)

	unreadLbl, err := gtk.LabelNew(strconv.Itoa(len(messages.Messages)))
	utils.HandleError(err)

	iconName := "go-next"

	seeMoreBtn, err := gtk.ButtonNewFromIconName(iconName, gtk.ICON_SIZE_BUTTON)
	utils.HandleError(err)

	seeMoreBtn.SetTooltipText("List messages")

	seeMoreBtn.Connect("clicked", func() {
		inbox.Show(box, messages)
	})

	openBtn, err := gtk.ButtonNewFromIconName("go-up", gtk.ICON_SIZE_BUTTON)
	utils.HandleError(err)

	openBtn.SetTooltipText("Open inbox on browser")

	openBtn.Connect("clicked", func() {
		//url := mail.GetWebURL()
		// TODO: cross platform?
		//_ = exec.Command("xdg-open", url).Start()
	})

	container.PackStart(emailLbl, false, false, 0)
	container.PackEnd(seeMoreBtn, false, false, 1)
	container.PackEnd(openBtn, false, false, 1)
	container.PackEnd(unreadLbl, false, false, 1)

	return container
}

func createInboxList() *gtk.ScrolledWindow {
	scroller, err := gtk.ScrolledWindowNew(nil, nil)
	utils.HandleError(err)

	container, err := gtk.GridNew()
	utils.HandleError(err)

	scroller.Add(container)

	container.SetRowSpacing(5)
	container.SetColumnHomogeneous(true)
	container.SetMarginStart(5)
	container.SetMarginEnd(5)

	spinner, err := gtk.SpinnerNew()
	utils.HandleError(err)

	spinner.Start()

	container.Attach(spinner, 0, 0, 1, 1)

	go func() {
		var messages []*types.CachedUnreadMessages

		inboxes, err := client.ListInboxes()
		if err != nil {
			panic(err)
		}

		for i := range inboxes {
			msgs, err := client.FetchUnreadMessagesIn(i)
			messages = append(messages, msgs)
			if err != nil {
				panic(err)
			}
		}

		glib.IdleAdd(func() {
			spinner.Destroy()
			for i, m := range inboxes {
				container.Attach(createInboxItem(m, messages[i]), 0, i, 1, 1)
			}
			container.ShowAll()
		})
	}()

	return scroller
}
