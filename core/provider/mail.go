package provider

var (
	Factories = map[string]MailProviderFactory{}
)

type MailProvider interface {
	FetchUnreadMessages() (ids []MailMessage, err error)
	FetchMessage(id string) (message MailMessage, err error)
	MarkMessageAsRead(id string) error
	GetAddress() string
	GetWebURL() string
}

type MailProviderFactory func(info map[string]interface{}) (MailProvider, error)
