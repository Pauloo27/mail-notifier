package provider

type MailMessage interface {
	GetID() string
	GetSubject() string
	GetBody() string
	GetSenderAddress() string
	GetRecipientAddress() string
}
