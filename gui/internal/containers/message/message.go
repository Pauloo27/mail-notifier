package message

import (
	"github.com/Pauloo27/mail-notifier/gui/utils"
	"github.com/Pauloo27/mail-notifier/socket/common/types"
	"github.com/gotk3/gotk3/gtk"
)

func createMessageHeaderPreview(msg *types.CachedMailMessage, contentType string) *gtk.Box {
	container, err := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 1)
	utils.HandleError(err)

	subjectLbl, err := gtk.LabelNew(utils.AddEllipsis(msg.Subject, 60))
	utils.HandleError(err)
	subjectLbl.SetHAlign(gtk.ALIGN_START)
	subjectLbl.SetTooltipText(msg.Subject)

	fromLbl, err := gtk.LabelNew(utils.AddEllipsis(msg.From, 60))
	utils.HandleError(err)
	fromLbl.SetHAlign(gtk.ALIGN_START)
	fromLbl.SetTooltipText(msg.From)

	date := msg.Date.Format("2006-01-02 03:04:05 PM")

	dateLbl, err := gtk.LabelNew(date)
	utils.HandleError(err)
	dateLbl.SetHAlign(gtk.ALIGN_START)
	dateLbl.SetTooltipText(date)

	container.PackStart(subjectLbl, false, false, 0)
	container.PackStart(fromLbl, false, false, 0)
	container.PackStart(dateLbl, false, false, 0)

	return container
}

func createMessageContentPreview(msg *types.CachedMailMessage, contentType string) *gtk.ScrolledWindow {
	scroller, err := gtk.ScrolledWindowNew(nil, nil)
	utils.HandleError(err)

	textBuf, err := gtk.TextBufferNew(nil)
	utils.HandleError(err)
	textBuf.SetText(string(msg.TextContents[contentType]))

	textView, err := gtk.TextViewNewWithBuffer(textBuf)
	utils.HandleError(err)
	textView.SetWrapMode(gtk.WRAP_WORD_CHAR)
	textView.SetEditable(false)

	scroller.Add(textView)
	scroller.SetVExpand(true)

	return scroller
}
