package gmail

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/Pauloo27/logger"
	"github.com/Pauloo27/mail-notifier/core/provider"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/gmail/v1"
)

var _ provider.MailBox = Gmail{}

type Gmail struct {
	Client      *http.Client
	Config      *oauth2.Config
	Service     *gmail.Service
	mailAddress string
	userID      int
	hideSpam    bool
}

var (
	mailAddressRegex = regexp.MustCompile(`(<(\S+@\S+)>)`)
)

func init() {
	provider.Factories["gmail"] = func(info map[string]interface{}) (provider.MailBox, error) {
		m, err := NewGmail(info["credentials"].(string), int(info["id"].(float64)), info["hide_spam"].(bool))
		if err != nil {
			return nil, err
		}
		token, err := m.ResolveTokenFromFile(info["token"].(string))
		if err != nil {
			return nil, err
		}
		err = m.LoginWithToken(token)
		return m, err
	}
}

func NewGmail(credentialsFilePath string, userID int, hideSpam bool) (*Gmail, error) {
	buf, err := os.ReadFile(credentialsFilePath)
	if err != nil {
		return nil, err
	}
	config, err := google.ConfigFromJSON(buf, gmail.GmailReadonlyScope)
	if err != nil {
		return nil, err
	}
	return &Gmail{
		hideSpam: hideSpam,
		userID:   userID,
		Config:   config,
	}, nil
}

func (m Gmail) FetchMessage(id string) (message provider.MailMessage, err error) {
	msg, err := m.Service.Users.Messages.Get("me", id).Do()
	if err != nil {
		return
	}

	contents := make(map[string][]byte)

	if len(msg.Payload.Parts) == 0 {
		contentType := msg.Payload.MimeType
		if strings.HasPrefix(contentType, "text/") {
			var b []byte
			b, err = base64.URLEncoding.DecodeString(msg.Payload.Body.Data)
			if err != nil {
				return
			}
			contents[contentType] = b
		}
	} else {
		for _, part := range msg.Payload.Parts {
			contentType := part.MimeType
			if strings.HasPrefix(contentType, "text/") {
				var b []byte
				b, err = base64.URLEncoding.DecodeString(part.Body.Data)
				if err != nil {
					logger.Error(err)
					return
				}
				contents[contentType] = b
			}
		}
	}

	var from, subject string
	var to []string
	for _, m := range msg.Payload.Headers {
		switch m.Name {
		case "Subject":
			subject = m.Value
		case "From":
			matches := mailAddressRegex.FindAllStringSubmatch(m.Value, 1)
			if len(matches) != 0 {
				from = matches[0][2]
			} else {
				from = m.Value
			}
		case "To":
			matches := mailAddressRegex.FindAllStringSubmatch(m.Value, -1)
			for _, match := range matches {
				to = append(to, match[2])
			}
			if len(to) == 0 {
				to = []string{m.Value}
			}
		}
	}
	message = GmailMessage{
		id:   id,
		mail: &m,
		data: &gmailMessageData{
			loaded:       true,
			subject:      subject,
			from:         from,
			textContents: contents,
			to:           to,
			date:         time.Unix(msg.InternalDate/1000, 0),
		},
	}
	return
}

func (m Gmail) GetAddress() string {
	return m.mailAddress
}

func (m Gmail) GetWebURL() string {
	return fmt.Sprintf("https://mail.google.com/mail/u/%d", m.userID)
}

func (m Gmail) MarkMessageAsRead(id string) (err error) {
	_, err = m.Service.Users.Messages.Modify("me", id, &gmail.ModifyMessageRequest{
		RemoveLabelIds: []string{"UNREAD"},
	}).Do()
	return err
}

func (m Gmail) FetchUnreadMessages() (messages []provider.MailMessage, err error) {
	query := m.Service.Users.Messages.List("me").IncludeSpamTrash(!m.hideSpam).LabelIds("UNREAD")

	res, err := query.Do()
	if err != nil {
		return nil, err
	}
	for _, message := range res.Messages {
		messages = append(messages, GmailMessage{
			mail: &m,
			id:   message.Id,
			data: &gmailMessageData{loaded: false},
		})
	}

	return messages, err
}
