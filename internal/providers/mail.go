package providers

type MailProvider interface {
	FetchMessages(maxMessages int, onlyUnread bool) (ids []string, count int, err error)
}
