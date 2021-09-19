package mail

import (
	"fmt"

	"github.com/Pauloo27/gmail-notifier/utils"
	"github.com/emersion/go-imap/client"
)

type Mail struct {
	Host, Username, Password string
	Port                     int
}

func (m *Mail) FetchMessages() uint32 {
	c, err := client.DialTLS(fmt.Sprintf("%s:%d", m.Host, m.Port), nil)
	utils.HandleFatal("Cannot connect to imap host", err)
	defer c.Logout()
	if err := c.Login(m.Username, m.Password); err != nil {
		utils.HandleFatal("Cannot login into imap host", err)
	}

	mbox, err := c.Select("INBOX", true)
	utils.HandleFatal("Cannot get INBOX mailbox", err)
	return mbox.Unseen
}
