package home

import (
	"github.com/Pauloo27/mail-notifier/gui/internal/config"
	"github.com/Pauloo27/mail-notifier/gui/utils"
	"github.com/Pauloo27/mail-notifier/internal/providers"
	"github.com/gotk3/gotk3/gtk"
)

func createInboxItem(email string, ok bool) *gtk.Box {
	container, err := gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 5)
	utils.HandleError(err)

	emailLbl, err := gtk.LabelNew(email)
	utils.HandleError(err)

	unreadLbl, err := gtk.LabelNew("10")
	utils.HandleError(err)

	var iconName string

	if ok {
		iconName = "go-next"
	} else {
		iconName = "security-low"
	}

	seeMoreBtn, err := gtk.ButtonNewFromIconName(iconName, gtk.ICON_SIZE_BUTTON)
	utils.HandleError(err)

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

	for i, provider := range config.Config.Providers {
		mail, err := providers.Factories[provider.Type](provider.Info)
		if err != nil {
			container.Attach(createInboxItem("invalid", false), 0, i, 1, 1)
			continue
		}
		container.Attach(createInboxItem(mail.GetAddress(), true), 0, i, 1, 1)
	}

	return scroller
}
