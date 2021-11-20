package provider

var (
	Factories = map[string]MailProviderFactory{}
)

type MailProvider interface {
	FetchMessages(onlyUnread bool) (ids []string, count int, err error)
	GetAddress() string
}

type MailProviderFactory func(info map[string]interface{}) (MailProvider, error)
