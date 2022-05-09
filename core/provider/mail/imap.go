package mail

import (
	"errors"
	"fmt"
	"strconv"
	"sync"
	"time"

	"github.com/Pauloo27/mail-notifier/core/provider"
	"github.com/emersion/go-imap"
	"github.com/emersion/go-imap/client"
	gomsg "github.com/emersion/go-message"
	"github.com/emersion/go-message/mail"
)

var _ provider.MailBox = Mail{}

type Mail struct {
	Host, Username, Password string
	Port                     int

	webURL string
	client *client.Client
	lock   *sync.Mutex
}

func init() {
	provider.Factories["imap"] = func(info map[string]interface{}) (provider.MailBox, error) {
		return NewMail(info["host"].(string), int(info["port"].(float64)), info["username"].(string), info["password"].(string), info["url"].(string))
	}
}

func (m *Mail) Connect() error {
	c, err := client.DialTLS(fmt.Sprintf("%s:%d", m.Host, m.Port), nil)
	m.client = c
	return err
}

func (m Mail) GetWebURL() string {
	return m.webURL
}

func (m Mail) GetAddress() string {
	return m.Username // FIXME: the username is always not the complete address,
	// maybe i can  get it from imap?
}

func NewMail(host string, port int, username, password string, webURL string) (Mail, error) {
	m := Mail{
		Host:     host,
		Port:     port,
		Username: username,
		Password: password,
		lock:     &sync.Mutex{},
		webURL:   webURL,
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

func (m Mail) MarkMessageAsRead(id string) (err error) {
	m.lock.Lock()
	defer m.lock.Unlock()

	_, err = m.client.Select("INBOX", false)
	if err != nil {
		return
	}
	defer m.client.Unselect()

	numericID, err := strconv.Atoi(id)
	if err != nil {
		return err
	}

	seqSet := new(imap.SeqSet)
	seqSet.AddNum(uint32(numericID))
	item := imap.FormatFlagsOp(imap.AddFlags, true)
	flags := []interface{}{imap.SeenFlag}
	err = m.client.Store(seqSet, item, flags, nil)

	return
}

func (m Mail) FetchMessage(id string) (message provider.MailMessage, err error) {
	m.lock.Lock()
	defer m.lock.Unlock()

	_, err = m.client.Select("INBOX", true)
	if err != nil {
		return
	}
	defer m.client.Unselect()

	seq := new(imap.SeqSet)
	seq.Add(id)

	msgCh := make(chan *imap.Message, 1)

	var section imap.BodySectionName
	items := []imap.FetchItem{section.FetchItem()}

	err = m.client.Fetch(seq, items, msgCh)
	if err != nil {
		return
	}

	msg := <-msgCh
	if msg == nil {
		err = errors.New("not found")
		return
	}

	body := msg.GetBody(&section)

	e, err := gomsg.Read(body)
	if err != nil {
		return
	}

	mr := mail.NewReader(e)

	var date time.Time
	var subject, from string
	var to []string
	var addrs []*mail.Address

	header := mr.Header
	date, err = header.Date()
	if err != nil {
		return
	}

	subject, err = header.Subject()
	if err != nil {
		return
	}

	addrs, err = header.AddressList("From")
	if err != nil {
		return
	}
	from = addrs[0].Address

	addrs, err = header.AddressList("To")
	if err != nil {
		return
	}

	for _, add := range addrs {
		to = append(to, add.Address)
	}

	message = MailMessage{
		id: id,
		data: &mailMessageData{
			date:    date,
			from:    from,
			to:      to,
			subject: subject,
			loaded:  true,
		},
	}

	return
}

func (m Mail) FetchUnreadMessages() (messages []provider.MailMessage, err error) {
	m.lock.Lock()
	defer m.lock.Unlock()

	criteria := imap.NewSearchCriteria()

	criteria.WithoutFlags = []string{imap.SeenFlag}

	_, err = m.client.Select("INBOX", true)
	if err != nil {
		return
	}

	defer m.client.Unselect()

	criteria.Since = time.Now().AddDate(-1, 0, 0) // limit search on 1 year

	var rawIDs []uint32
	rawIDs, err = m.client.Search(criteria)

	for _, id := range rawIDs {
		messages = append(messages, MailMessage{
			id:   strconv.Itoa(int(id)),
			data: &mailMessageData{loaded: false},
			mail: &m,
		})
	}

	return
}
