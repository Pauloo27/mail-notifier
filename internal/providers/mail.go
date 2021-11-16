package providers

var (
	Providers = map[string]MailProviderFactory{}
)

type MailProvider interface {
	FetchMessages(maxMessages int, onlyUnread bool) (ids []string, count int, err error)
}

type MailProviderFactory func(info map[string]interface{}) (MailProvider, error)
