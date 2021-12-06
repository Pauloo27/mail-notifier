package inbox

import (
	"os/exec"

	"github.com/Pauloo27/mail-notifier/gui/utils"
	"github.com/Pauloo27/mail-notifier/socket/common/types"
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
)

var (
	messages      *types.CachedUnreadMessages
	mainContainer *gtk.Grid
)

func asyncLoad(container *gtk.Grid, box *types.Inbox, messages *types.CachedUnreadMessages) {
	spinner, err := gtk.SpinnerNew()
	utils.HandleError(err)

	spinner.Start()

	container.Attach(spinner, 0, 0, 1, 1)

	go func() {
		glib.IdleAdd(func() {
			children := container.GetChildren()

			children.Foreach(func(item interface{}) {
				wid := item.(*gtk.Widget)
				wid.Destroy()
			})

			for i, message := range messages.Messages {
				container.Attach(createMessageItem(box, message), 0, i, 1, 1)
			}
			container.ShowAll()
		})
	}()

}

func createMessageItem(box *types.Inbox, message *types.CachedMailMessage) *gtk.Box {
	container, err := gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 5)
	utils.HandleError(err)

	leftContainer, err := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 5)
	utils.HandleError(err)

	subjectLbl, err := gtk.LabelNew(utils.AddEllipsis(message.Subject, 40))
	utils.HandleError(err)
	subjectLbl.SetHAlign(gtk.ALIGN_START)

	fromLbl, err := gtk.LabelNew(utils.AddEllipsis(message.From, 30))
	utils.HandleError(err)
	fromLbl.SetHAlign(gtk.ALIGN_START)

	leftContainer.PackStart(subjectLbl, false, false, 1)
	leftContainer.PackStart(fromLbl, false, false, 1)

	markAsReadBtn, err := gtk.ButtonNewFromIconName("mail-read", gtk.ICON_SIZE_BUTTON)
	utils.HandleError(err)

	markAsReadBtn.SetVAlign(gtk.ALIGN_CENTER)

	markAsReadBtn.Connect("clicked", func() {
		var newMessages []*types.CachedMailMessage
		for _, m := range messages.Messages {
			if m.ID != message.ID {
				newMessages = append(newMessages, m)
			}
		}
		messages.Messages = newMessages
		asyncLoad(mainContainer, box, messages)
		go func() {
			/* FIXME
			err := box.MarkMessageAsRead(message.GetID())
			if err != nil {
				panic(err)
			}
			*/
		}()
	})

	markAsReadBtn.SetTooltipText("Mark as read")

	openBtn, err := gtk.ButtonNewFromIconName("go-up", gtk.ICON_SIZE_BUTTON)
	utils.HandleError(err)

	openBtn.SetTooltipText("Open inbox on browser")
	openBtn.SetVAlign(gtk.ALIGN_CENTER)

	openBtn.Connect("clicked", func() {
		//url := mail.GetWebURL()
		url := "" // FIXME
		// TODO: cross platform?
		_ = exec.Command("xdg-open", url).Start()
	})

	container.PackStart(leftContainer, false, false, 1)
	container.PackEnd(markAsReadBtn, false, false, 1)
	container.PackEnd(openBtn, false, false, 1)

	return container
}

func createMessageList(box *types.Inbox, messagesParam *types.CachedUnreadMessages) *gtk.ScrolledWindow {
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

	asyncLoad(mainContainer, box, messages)

	return scroller
}
