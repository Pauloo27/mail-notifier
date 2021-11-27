package gmail

import (
	"context"
	"encoding/json"
	"os"

	"golang.org/x/oauth2"
	"google.golang.org/api/gmail/v1"
)

func (m *Gmail) GetLoginURL() string {
	return m.Config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
}

func (m *Gmail) ResolveToken(authCode string) (*oauth2.Token, error) {
	return m.Config.Exchange(context.TODO(), authCode)
}

func (m *Gmail) ResolveTokenFromFile(tokenFilePath string) (*oauth2.Token, error) {
	buf, err := os.ReadFile(tokenFilePath)
	if err != nil {
		return nil, err
	}
	var token oauth2.Token
	err = json.Unmarshal(buf, &token)
	return &token, err
}

func (m *Gmail) SaveTokenToFile(token *oauth2.Token, tokenFilePath string) error {
	buf, err := json.Marshal(token)
	if err != nil {
		return err
	}
	return os.WriteFile(tokenFilePath, buf, 0600)
}

func (m *Gmail) LoginWithToken(token *oauth2.Token) error {
	m.Client = m.Config.Client(context.Background(), token)
	service, err := gmail.New(m.Client)
	if err != nil {
		return err
	}
	m.Service = service
	p, err := service.Users.GetProfile("me").Do()
	if err != nil {
		return err
	}
	m.mailAddress = p.EmailAddress
	return nil
}
