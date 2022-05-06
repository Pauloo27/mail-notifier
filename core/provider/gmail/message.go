package gmail

import (
	"time"

	"github.com/Pauloo27/logger"
	"github.com/Pauloo27/mail-notifier/core/provider"
)

var _ provider.MailMessage = GmailMessage{}

type gmailMessageData struct {
	loaded        bool
	date          time.Time
	from, subject string
	textContents  map[string][]byte
	to            []string
}

type GmailMessage struct {
	mail *Gmail
	id   string

	data *gmailMessageData
}

func (m GmailMessage) GetID() string {
	return m.id
}

func (m *GmailMessage) load() error {
	if m.data.loaded {
		return nil
	}
	fullMsg, err := m.mail.FetchMessage(m.id)
	if err != nil {
		return err
	}
	*m.data = *fullMsg.(GmailMessage).data
	return nil
}

func (m GmailMessage) GetTextContents() map[string][]byte {
	if !m.data.loaded {
		if err := m.load(); err != nil {
			logger.Fatal(err)
		}
	}
	return m.data.textContents
}

func (m GmailMessage) GetSubject() string {
	if !m.data.loaded {
		if err := m.load(); err != nil {
			logger.Fatal(err)
		}
	}
	return m.data.subject
}

func (m GmailMessage) GetFrom() string {
	if !m.data.loaded {
		if err := m.load(); err != nil {
			logger.Fatal(err)
		}
	}
	return m.data.from
}

func (m GmailMessage) GetTo() []string {
	if !m.data.loaded {
		if err := m.load(); err != nil {
			logger.Fatal(err)
		}
	}
	return m.data.to
}

func (m GmailMessage) GetDate() time.Time {
	if !m.data.loaded {
		if err := m.load(); err != nil {
			logger.Fatal(err)
		}
	}
	return m.data.date
}
