package mail

import (
	"fmt"

	"github.com/emersion/go-imap"
	"github.com/emersion/go-imap/client"
)

type Mail struct {
	Host, Username, Password string
	Port                     int

	client *client.Client
}

func (m *Mail) Connect() error {
	c, err := client.DialTLS(fmt.Sprintf("%s:%d", m.Host, m.Port), nil)
	m.client = c
	return err
}

func NewMail(host string, port int, username, password string) (*Mail, error) {
	m := &Mail{
		Host:     host,
		Port:     port,
		Username: username,
		Password: password,
	}
	err := m.Connect()
	if err == nil {
		err = m.client.Login(username, password)
	}
	return m, err
}

func (m *Mail) Disconnect() error {
	return m.client.Logout()
}

func (m *Mail) FetchMessages(maxMessages uint32, onlyUnread bool) (ids []uint32, count int, err error) {
	criteria := imap.NewSearchCriteria()

	if onlyUnread {
		criteria.WithoutFlags = []string{imap.SeenFlag}
	}

	var box *imap.MailboxStatus

	box, err = m.client.Select("INBOX", true)
	if err != nil {
		return
	}

	if box.Messages < maxMessages {
		maxMessages = box.Messages
	}

	seq := new(imap.SeqSet)
	seq.AddRange(1, maxMessages)
	criteria.SeqNum = seq

	ids, err = m.client.Search(criteria)

	count = len(ids)

	return
}
