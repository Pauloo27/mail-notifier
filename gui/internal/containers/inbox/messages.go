package inbox

import (
	"github.com/Pauloo27/mail-notifier/gui/utils"
	"github.com/Pauloo27/mail-notifier/internal/provider"
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
)

var (
	messages      []provider.MailMessage
	mainContainer *gtk.Grid
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
			children := container.GetChildren()

			children.Foreach(func(item interface{}) {
				wid := item.(*gtk.Widget)
				wid.Destroy()
			})

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

	subjectLbl, err := gtk.LabelNew(utils.AddEllipsis(message.GetSubject(), 50))
	utils.HandleError(err)
	subjectLbl.SetHAlign(gtk.ALIGN_START)

	fromLbl, err := gtk.LabelNew(utils.AddEllipsis(message.GetFrom(), 40))
	utils.HandleError(err)
	fromLbl.SetHAlign(gtk.ALIGN_START)

	leftContainer.PackStart(subjectLbl, false, false, 1)
	leftContainer.PackStart(fromLbl, false, false, 1)

	markAsReadBtn, err := gtk.ButtonNewFromIconName("mail-read", gtk.ICON_SIZE_BUTTON)
	utils.HandleError(err)

	markAsReadBtn.SetVAlign(gtk.ALIGN_CENTER)

	markAsReadBtn.Connect("clicked", func() {
		var newMessages []provider.MailMessage
		for _, m := range messages {
			if m.GetID() != message.GetID() {
				newMessages = append(newMessages, m)
			}
		}
		messages = newMessages
		asyncLoad(mainContainer, mail, messages)
		go func() {
			err := mail.MarkMessageAsRead(message.GetID())
			if err != nil {
				panic(err) // FIXME
			}
		}()
	})

	markAsReadBtn.SetTooltipText("Mark as read")

	container.PackStart(leftContainer, false, false, 1)
	container.PackEnd(markAsReadBtn, false, false, 10)

	return container
}

func createMessageList(mail provider.MailProvider, messagesParam []provider.MailMessage) *gtk.ScrolledWindow {
	messages = messagesParam

	scroller, err := gtk.ScrolledWindowNew(nil, nil)
	utils.HandleError(err)

	mainContainer, err = gtk.GridNew()
	utils.HandleError(err)

	scroller.Add(mainContainer)

	mainContainer.SetRowSpacing(5)
	mainContainer.SetColumnHomogeneous(true)
	mainContainer.SetMarginStart(5)
	mainContainer.SetMarginEnd(5)

	asyncLoad(mainContainer, mail, messages)

	return scroller
}
