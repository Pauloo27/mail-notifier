package mail

import (
	"fmt"
	"strconv"
	"time"

	"github.com/Pauloo27/mail-notifier/internal/provider"
	"github.com/emersion/go-imap"
	"github.com/emersion/go-imap/client"
)

var _ provider.MailProvider = Mail{}

type Mail struct {
	Host, Username, Password string
	Port                     int

	client *client.Client
}

func init() {
	provider.Factories["imap"] = func(info map[string]interface{}) (provider.MailProvider, error) {
		return NewMail(info["host"].(string), int(info["port"].(float64)), info["username"].(string), info["password"].(string))
	}
}

func (m *Mail) Connect() error {
	c, err := client.DialTLS(fmt.Sprintf("%s:%d", m.Host, m.Port), nil)
	m.client = c
	return err
}

func (m Mail) GetAddress() string {
	return m.Username // FIXME
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

func (m *Mail) Disconnect() error {
	return m.client.Logout()
}

func (m Mail) FetchMessages(onlyUnread bool) (ids []string, count int, err error) {
	criteria := imap.NewSearchCriteria()

	if onlyUnread {
		criteria.WithoutFlags = []string{imap.SeenFlag}
	}

	_, err = m.client.Select("INBOX", true)
	if err != nil {
		return
	}

	criteria.Since = time.Now().AddDate(-1, 0, 0) // limit search on 1 year

	var rawIDs []uint32
	rawIDs, err = m.client.Search(criteria)

	for _, id := range rawIDs {
		ids = append(ids, strconv.Itoa(int(id)))
	}

	count = len(ids)

	return
}
