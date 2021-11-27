package inbox

import (
	"github.com/Pauloo27/mail-notifier/gui/utils"
	"github.com/Pauloo27/mail-notifier/internal/provider"
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
)

func asyncLoad(container *gtk.Grid, mail provider.MailProvider, messages []provider.MailMessage) {
	spinner, err := gtk.SpinnerNew()
	utils.HandleError(err)

	spinner.Start()

	container.Attach(spinner, 0, 0, 1, 1)

	go func() {
		for _, message := range messages {
			_ = message.GetFrom() // Just to trigger the lazy load
		}
		glib.IdleAdd(func() {
			spinner.Destroy()
			for i, message := range messages {
				container.Attach(createMessageItem(mail, message), 0, i, 1, 1)
			}
			container.ShowAll()
		})
	}()

}

func createMessageItem(mail provider.MailProvider, message provider.MailMessage) *gtk.Box {
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

	markAsReadBtn.SetVAlign(gtk.ALIGN_CENTER)

	markAsReadBtn.Connect("clicked", func() {
		err := mail.MarkMessageAsRead(message.GetID())
		if err != nil {
			panic(err) // FIXME
		}
	})

	container.PackStart(leftContainer, false, false, 1)
	container.PackEnd(markAsReadBtn, false, false, 10)

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

	asyncLoad(container, mail, messages)

	return scroller
}
