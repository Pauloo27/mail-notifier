package mail

import (
	"time"

	"github.com/Pauloo27/mail-notifier/internal/provider"
)

var _ provider.MailMessage = MailMessage{}

type MailMessage struct {
	id   string
	mail *Mail

	loaded        bool
	date          time.Time
	from, subject string
	to            []string
}

func (m *MailMessage) load() error {
	if m.loaded {
		return nil
	}
	fullMsg, err := m.mail.FetchMessage(m.id)
	if err != nil {
		return err
	}
	m.from = fullMsg.GetFrom()
	m.to = fullMsg.GetTo()
	m.subject = fullMsg.GetSubject()
	m.date = fullMsg.GetDate()
	return nil
}

func (m MailMessage) GetID() string {
	return m.id
}

func (m MailMessage) GetSubject() string {
	if !m.loaded {
		if err := m.load(); err != nil {
			panic(err)
		}
	}
	return m.subject
}

func (m MailMessage) GetFrom() string {
	if !m.loaded {
		if err := m.load(); err != nil {
			panic(err)
		}
	}
	return m.from
}

func (m MailMessage) GetTo() []string {
	if !m.loaded {
		if err := m.load(); err != nil {
			panic(err)
		}
	}
	return m.to
}

func (m MailMessage) GetDate() time.Time {
	if !m.loaded {
		if err := m.load(); err != nil {
			panic(err)
		}
	}
	return m.date
}
