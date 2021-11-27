package mail

import (
	"time"

	"github.com/Pauloo27/mail-notifier/internal/provider"
)

var _ provider.MailMessage = MailMessage{}

type mailMessageData struct {
	loaded        bool
	date          time.Time
	from, subject string
	to            []string
}

type MailMessage struct {
	id   string
	mail *Mail

	data *mailMessageData
}

func (m *MailMessage) load() error {
	if m.data.loaded {
		return nil
	}
	fullMsg, err := m.mail.FetchMessage(m.id)
	if err != nil {
		return err
	}
	*m.data = *fullMsg.(MailMessage).data
	return nil
}

func (m MailMessage) GetID() string {
	return m.id
}

func (m MailMessage) GetSubject() string {
	if !m.data.loaded {
		if err := m.load(); err != nil {
			panic(err)
		}
	}
	return m.data.subject
}

func (m MailMessage) GetFrom() string {
	if !m.data.loaded {
		if err := m.load(); err != nil {
			panic(err)
		}
	}
	return m.data.from
}

func (m MailMessage) GetTo() []string {
	if !m.data.loaded {
		if err := m.load(); err != nil {
			panic(err)
		}
	}
	return m.data.to
}

func (m MailMessage) GetDate() time.Time {
	if !m.data.loaded {
		if err := m.load(); err != nil {
			panic(err)
		}
	}
	return m.data.date
}
