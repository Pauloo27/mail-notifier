package gmail

import (
	"net/http"

	"golang.org/x/oauth2"
	"google.golang.org/api/gmail/v1"
)

type Gmail struct {
	Client  *http.Client
	Config  *oauth2.Config
	Service *gmail.Service
}

func (m *Gmail) FetchMessages(onlyUnread bool) (ids []string, count int, err error) {
	query := m.Service.Users.Messages.List("me").IncludeSpamTrash(true)
	if onlyUnread {
		query.LabelIds("UNREAD")
	}

	res, err := query.Do()
	if err != nil {
		return nil, 0, err
	}

	for _, message := range res.Messages {
		ids = append(ids, message.Id)
	}

	return ids, len(ids), err
}
