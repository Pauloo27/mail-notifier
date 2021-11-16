package home

import (
	"github.com/Pauloo27/mail-notifier/gui/utils"
	"github.com/gotk3/gotk3/gtk"
)

func createInboxItem(email string) *gtk.Box {
	container, err := gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 5)
	utils.HandleError(err)

	emailLbl, err := gtk.LabelNew(email)
	utils.HandleError(err)

	unreadLbl, err := gtk.LabelNew("10")
	utils.HandleError(err)

	seeMoreBtn, err := gtk.ButtonNewFromIconName("go-next", gtk.ICON_SIZE_BUTTON)
	utils.HandleError(err)

	container.PackStart(emailLbl, false, false, 0)
	container.PackEnd(seeMoreBtn, false, false, 10)
	container.PackEnd(unreadLbl, false, false, 1)

	return container
}

var dummyInboxes = []string{
	"test@example.com",
	"test@example.com",
	"test@example.com",
	"test@example.com",
	"test@example.com",
	"test@example.com",
	"test@example.com",
	"test@example.com",
	"test@example.com",
	"test@example.com",
	"test@example.com",
	"test@example.com",
	"test@example.com",
	"test@example.com",
	"test@example.com",
	"test@example.com",
	"test@example.com",
	"test@example.com",
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

	for i, inbox := range dummyInboxes {
		container.Attach(createInboxItem(inbox), 0, i, 1, 1)
	}

	return scroller
}
