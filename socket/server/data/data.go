package data

import (
	"errors"
	"time"

	"github.com/Pauloo27/mail-notifier/core/provider"
	"github.com/Pauloo27/mail-notifier/core/storage"
	"github.com/Pauloo27/mail-notifier/socket/common/types"
)

var (
	config         *storage.Config
	inboxes        []provider.MailBox
	inboxMessages  = make(map[int]map[string]*types.CachedMailMessage)
	unreadMessages = make(map[int]*types.CachedUnreadMessages)

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
		inboxMessages[i] = make(map[string]*types.CachedMailMessage)
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
	inboxMessages[inboxID][messageID] = &types.CachedMailMessage{
		MailMessage: msg,
		FechedAt:    time.Now(),
	}
	return nil
}

func GetMessage(inboxID int, messageID string) (*types.CachedMailMessage, error) {
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

func fetchUnreadMessage(inboxID int) error {
	if config == nil {
		return ErrConfigNotLoaded
	}
	if inboxID == len(inboxes) {
		return ErrInvalidInbox
	}
	inbox := inboxes[inboxID]

	msgs, err := inbox.FetchUnreadMessages()
	if err != nil {
		return err
	}

	var unreadMsgs []*types.CachedMailMessage

	for _, msg := range msgs {
		unreadMsg, err := GetMessage(inboxID, msg.GetID())
		if err != nil {
			return err
		}
		unreadMsgs = append(unreadMsgs, unreadMsg)
	}

	unreadMessages[inboxID] = &types.CachedUnreadMessages{
		Messages: unreadMsgs,
		FechedAt: time.Now(),
	}

	return nil
}

func GetAllUnreadMessages() ([]*types.CachedUnreadMessages, error) {
	var messages []*types.CachedUnreadMessages
	for i := range inboxes {
		msgs, err := GetUnreadMessagesIn(i)
		if err != nil {
			return nil, err
		}
		messages = append(messages, msgs)
	}
	return messages, nil
}

func GetUnreadMessagesIn(inboxID int) (*types.CachedUnreadMessages, error) {
	unreadMessages, found := unreadMessages[inboxID]
	if !found {
		err := fetchUnreadMessage(inboxID)
		if err != nil {
			return nil, err
		}
		return GetUnreadMessagesIn(inboxID)
	}
	return unreadMessages, nil
}