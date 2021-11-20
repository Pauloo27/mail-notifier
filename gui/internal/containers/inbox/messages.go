package inbox

import (
	"github.com/Pauloo27/mail-notifier/gui/utils"
	"github.com/Pauloo27/mail-notifier/internal/provider"
	"github.com/gotk3/gotk3/gtk"
)

func createMessageItem(message provider.MailMessage) *gtk.Box {
	container, err := gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 5)
	utils.HandleError(err)

	unreadLbl, err := gtk.LabelNew(message.GetID())
	utils.HandleError(err)

	container.PackStart(unreadLbl, false, false, 1)

	return container
}

func createMessageList(mail provider.MailProvider, messages []provider.MailMessage) *gtk.ScrolledWindow {
	scroller, err := gtk.ScrolledWindowNew(nil, nil)
	utils.HandleError(err)

	container, err := gtk.GridNew()
	utils.HandleError(err)

	scroller.Add(container)

	container.SetRowSpacing(5)
	container.SetColumnHomogeneous(true)
	container.SetMarginStart(5)
	container.SetMarginEnd(5)

	for i, message := range messages {
		container.Attach(createMessageItem(message), 0, i, 1, 1)
	}

	return scroller
}
