package provider

var (
	Factories = map[string]MailBoxFactory{}
)

type MailBox interface {
	FetchUnreadMessages() (ids []MailMessage, err error)
	FetchMessage(id string) (message MailMessage, err error)
	MarkMessageAsRead(id string) error
	GetAddress() string
	GetWebURL() string
}

type MailBoxFactory func(info map[string]interface{}) (MailBox, error)
