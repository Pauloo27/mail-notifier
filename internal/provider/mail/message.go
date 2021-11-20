package mail

import "github.com/Pauloo27/mail-notifier/internal/provider"

var _ provider.MailMessage = MailMessage{}

type MailMessage struct {
	id   string
	mail *Mail
}

func (m MailMessage) GetID() string {
	return m.id
}

func (m MailMessage) GetSubject() string {
	// TODO: lady load info
	return ""
}

func (m MailMessage) GetBody() string {
	return ""
}

func (m MailMessage) GetSenderAddress() string {
	return ""
}

func (m MailMessage) GetRecipientAddress() string {
	return ""
}
