package mail

import (
	"fmt"
	"strconv"

	"github.com/Pauloo27/mail-notifier/internal/providers"

	"github.com/emersion/go-imap"
	"github.com/emersion/go-imap/client"
)

type Mail struct {
	Host, Username, Password string
	Port                     int

	client *client.Client
}

var _ providers.MailProvider = Mail{}

func (m Mail) Connect() error {
	c, err := client.DialTLS(fmt.Sprintf("%s:%d", m.Host, m.Port), nil)
	m.client = c
	return err
}

func NewMail(host string, port int, username, password string) (Mail, error) {
	m := Mail{
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

func (m Mail) Disconnect() error {
	return m.client.Logout()
}

func (m Mail) FetchMessages(maxMessages int, onlyUnread bool) (ids []string, count int, err error) {
	criteria := imap.NewSearchCriteria()

	if onlyUnread {
		criteria.WithoutFlags = []string{imap.SeenFlag}
	}

	var box *imap.MailboxStatus

	box, err = m.client.Select("INBOX", true)
	if err != nil {
		return
	}

	maxMessagesUint := uint32(maxMessages)

	if box.Messages < maxMessagesUint {
		maxMessagesUint = box.Messages
	}

	seq := new(imap.SeqSet)
	seq.AddRange(1, maxMessagesUint)
	criteria.SeqNum = seq

	var rawIDs []uint32
	rawIDs, err = m.client.Search(criteria)

	for _, id := range rawIDs {
		ids = append(ids, strconv.Itoa(int(id)))
	}

	count = len(ids)

	return
}
