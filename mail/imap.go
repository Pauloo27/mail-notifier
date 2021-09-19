package mail

import (
	"fmt"

	"github.com/emersion/go-imap/client"
)

type Mail struct {
	Host, Username, Password string
	Port                     int
}

func (m *Mail) FetchMessages() uint32 {
	c, err := client.DialTLS(fmt.Sprintf("%s:%d", m.Host, m.Port), nil)
	if err != nil {
		panic(err)
	}
	defer c.Logout()
	if err := c.Login(m.Username, m.Password); err != nil {
		panic(err)
	}

	mbox, err := c.Select("INBOX", true)
	if err != nil {
		panic(err)
	}
	return mbox.Unseen
}
