package data

import (
	"errors"

	"github.com/Pauloo27/mail-notifier/core/provider"
	"github.com/Pauloo27/mail-notifier/core/storage"
)

var (
	config        *storage.Config
	inboxes       []provider.MailBox
	inboxMessages = make(map[int]map[string]*provider.MailMessage)

	ErrConfigNotLoaded = errors.New("config not loaded")
	ErrInvalidInbox    = errors.New("invalid inbox")
)

func LoadConfig() (err error) {
	config, err = storage.LoadConfig()
	return
}

func ConnectToInboxes() (err error) {
	if config == nil {
		return ErrConfigNotLoaded
	}
	for i, p := range config.Providers {
		inbox, err := provider.Factories[p.Type](p.Info)
		if err != nil {
			return err
		}
		inboxes = append(inboxes, inbox)
		inboxMessages[i] = make(map[string]*provider.MailMessage)
	}
	return nil
}

func GetInboxes() ([]provider.MailBox, error) {
	if config == nil {
		return nil, ErrConfigNotLoaded
	}
	return inboxes, nil
}

func fetchMessage(inboxID int, messageID string) error {
	if config == nil {
		return ErrConfigNotLoaded
	}
	if inboxID == len(inboxes) {
		return ErrInvalidInbox
	}
	inbox := inboxes[inboxID]
	msg, err := inbox.FetchMessage(messageID)
	if err != nil {
		return err
	}
	inboxMessages[inboxID][messageID] = &msg
	return nil
}

func GetMessage(inboxID int, messageID string) (*provider.MailMessage, error) {
	inbox, found := inboxMessages[inboxID]
	if !found {
		return nil, ErrInvalidInbox
	}
	message, found := inbox[messageID]
	if !found {
		err := fetchMessage(inboxID, messageID)
		if err != nil {
			return nil, err
		}
		return GetMessage(inboxID, messageID)
	}
	return message, nil
}
