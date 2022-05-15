package inbox

import (
	"os/exec"
	"strings"

	"github.com/Pauloo27/logger"
	"github.com/Pauloo27/mail-notifier/gui/internal/containers/message"
	"github.com/Pauloo27/mail-notifier/gui/utils"
	"github.com/Pauloo27/mail-notifier/socket/client"
	"github.com/Pauloo27/mail-notifier/socket/common/types"
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
)

var (
	messages      *types.CachedUnreadMessages
	mainContainer *gtk.Grid
)

func asyncLoad(container *gtk.Grid, c *client.Client, box *types.Inbox, messages *types.CachedUnreadMessages) {
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
				container.Attach(createMessageItem(c, box, message), 0, i, 1, 1)
			}
			container.ShowAll()
		})
	}()

}

func createMessageItem(c *client.Client, box *types.Inbox, msg *types.CachedMailMessage) *gtk.Box {
	container, err := gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 5)
	utils.HandleError(err)

	leftContainer, err := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 1)
	utils.HandleError(err)

	subjectLbl, err := gtk.LabelNew(utils.AddEllipsis(msg.Subject, 40))
	utils.HandleError(err)
	subjectLbl.SetHAlign(gtk.ALIGN_START)
	subjectLbl.SetTooltipText(msg.Subject)

	fromLbl, err := gtk.LabelNew(utils.AddEllipsis(msg.From, 30))
	utils.HandleError(err)
	fromLbl.SetHAlign(gtk.ALIGN_START)
	fromLbl.SetTooltipText(msg.From)

	date := msg.Date.Format("2006-01-02 03:04:05 PM")

	dateLbl, err := gtk.LabelNew(date)
	utils.HandleError(err)
	dateLbl.SetHAlign(gtk.ALIGN_START)
	dateLbl.SetTooltipText(date)

	leftContainer.PackStart(subjectLbl, false, false, 0)
	leftContainer.PackStart(fromLbl, false, false, 0)
	leftContainer.PackStart(dateLbl, false, false, 0)

	var textContentType string
	for contentType := range msg.TextContents {
		if strings.HasPrefix(contentType, "text/plain") {
			textContentType = contentType
			break
		}
	}

	previewMailBtn, err := gtk.ButtonNewFromIconName("go-next", gtk.ICON_SIZE_BUTTON)
	utils.HandleError(err)
	previewMailBtn.SetVAlign(gtk.ALIGN_CENTER)
	previewMailBtn.SetTooltipText("Preview mail content")
	previewMailBtn.SetSensitive(textContentType != "")
	previewMailBtn.Connect("clicked", func() {
		if textContentType != "" {
			message.Show(msg, textContentType)
		}
	})

	markAsReadBtn, err := gtk.ButtonNewFromIconName("mail-read", gtk.ICON_SIZE_BUTTON)
	utils.HandleError(err)

	markAsReadBtn.SetVAlign(gtk.ALIGN_CENTER)

	markAsReadBtn.Connect("clicked", func() {
		var newMessages []*types.CachedMailMessage
		for _, m := range messages.Messages {
			if m.ID != msg.ID {
				newMessages = append(newMessages, m)
			}
		}
		messages.Messages = newMessages
		asyncLoad(mainContainer, c, box, messages)
		go func() {
			err := c.MarkMessageAsRead(box.ID, msg.ID)
			if err != nil {
				logger.Fatal(err)
			}
		}()
	})

	markAsReadBtn.SetTooltipText("Mark as read")

	openBrowserBtn, err := gtk.ButtonNewFromIconName("go-up", gtk.ICON_SIZE_BUTTON)
	utils.HandleError(err)

	openBrowserBtn.SetTooltipText("Open inbox on browser")
	openBrowserBtn.SetVAlign(gtk.ALIGN_CENTER)

	openBrowserBtn.Connect("clicked", func() {
		url := box.WebURL
		// TODO: cross platform?
		_ = exec.Command("xdg-open", url).Start()
	})

	container.PackStart(leftContainer, false, false, 1)
	container.PackEnd(previewMailBtn, false, false, 1)
	container.PackEnd(markAsReadBtn, false, false, 1)
	container.PackEnd(openBrowserBtn, false, false, 1)

	return container
}

func createMessageList(c *client.Client, box *types.Inbox, messagesParam *types.CachedUnreadMessages) *gtk.ScrolledWindow {
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

	asyncLoad(mainContainer, c, box, messages)

	return scroller
}
