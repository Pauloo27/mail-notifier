package inbox

import (
	"github.com/Pauloo27/mail-notifier/gui/utils"
	"github.com/Pauloo27/mail-notifier/internal/provider"
	"github.com/gotk3/gotk3/gtk"
)

func createMessageItem(message provider.MailMessage) *gtk.Box {
	container, err := gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 5)
	utils.HandleError(err)

	leftContainer, err := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 5)
	utils.HandleError(err)

	subjectLbl, err := gtk.LabelNew(message.GetSubject())
	utils.HandleError(err)
	subjectLbl.SetHAlign(gtk.ALIGN_START)

	fromLbl, err := gtk.LabelNew(message.GetFrom())
	utils.HandleError(err)
	fromLbl.SetHAlign(gtk.ALIGN_START)

	leftContainer.PackStart(subjectLbl, false, false, 1)
	leftContainer.PackStart(fromLbl, false, false, 1)

	markAsReadBtn, err := gtk.ButtonNewFromIconName("mail-read", gtk.ICON_SIZE_BUTTON)
	utils.HandleError(err)

	markAsReadBtn.Connect("clicked", func() {
		// TODO: write me please
	})

	container.PackStart(leftContainer, false, false, 1)
	container.PackEnd(markAsReadBtn, false, true, 1)

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
