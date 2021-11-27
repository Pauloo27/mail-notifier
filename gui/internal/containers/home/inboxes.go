package home

import (
	"strconv"

	"github.com/Pauloo27/mail-notifier/gui/internal/config"
	"github.com/Pauloo27/mail-notifier/gui/internal/containers/inbox"
	"github.com/Pauloo27/mail-notifier/gui/utils"
	"github.com/Pauloo27/mail-notifier/internal/provider"
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
)

func createInboxItem(mail provider.MailProvider, messages []provider.MailMessage) *gtk.Box {
	container, err := gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 5)
	utils.HandleError(err)

	ok := mail != nil

	address := "invalid"
	if ok {
		address = mail.GetAddress()
	}

	emailLbl, err := gtk.LabelNew(address)
	utils.HandleError(err)

	unreadLbl, err := gtk.LabelNew(strconv.Itoa(len(messages)))
	utils.HandleError(err)

	var iconName string

	if ok {
		iconName = "go-next"
	} else {
		iconName = "security-low"
	}

	seeMoreBtn, err := gtk.ButtonNewFromIconName(iconName, gtk.ICON_SIZE_BUTTON)
	utils.HandleError(err)

	seeMoreBtn.SetTooltipText("List messages")

	seeMoreBtn.Connect("clicked", func() {
		inbox.Show(mail, messages)
	})

	container.PackStart(emailLbl, false, false, 0)
	container.PackEnd(seeMoreBtn, false, false, 10)
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
		var loadedMails []provider.MailProvider
		var messages [][]provider.MailMessage

		for _, p := range config.Config.Providers {
			mail, err := provider.Factories[p.Type](p.Info)
			if err == nil {
				mailMessages, err := mail.FetchMessages(true)
				if err == nil {
					loadedMails = append(loadedMails, mail)
					messages = append(messages, mailMessages)
					continue
				}
			}
			loadedMails = append(loadedMails, nil)
			messages = append(messages, nil)
		}

		glib.IdleAdd(func() {
			spinner.Destroy()
			for i, m := range loadedMails {
				container.Attach(createInboxItem(m, messages[i]), 0, i, 1, 1)
			}
			container.ShowAll()
		})
	}()

	return scroller
}
