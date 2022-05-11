package message

import (
	"github.com/Pauloo27/mail-notifier/gui/utils"
	"github.com/Pauloo27/mail-notifier/socket/common/types"
	"github.com/gotk3/gotk3/gtk"
)

func createMessagePreview(msg *types.CachedMailMessage, contentType string) *gtk.ScrolledWindow {
	scroller, err := gtk.ScrolledWindowNew(nil, nil)
	utils.HandleError(err)

	textBuf, err := gtk.TextBufferNew(nil)
	utils.HandleError(err)
	textBuf.SetText(string(msg.TextContents[contentType]))

	textView, err := gtk.TextViewNewWithBuffer(textBuf)
	utils.HandleError(err)
	textView.SetWrapMode(gtk.WRAP_WORD_CHAR)

	scroller.Add(textView)
	scroller.SetVExpand(true)

	return scroller
}
