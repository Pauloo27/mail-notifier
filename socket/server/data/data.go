package data

import (
	"errors"

	"github.com/Pauloo27/mail-notifier/core/provider"
	"github.com/Pauloo27/mail-notifier/core/storage"
)

var (
	config  *storage.Config
	inboxes []provider.MailBox

	ErrConfigNotLoaded = errors.New("config not loaded")
)

func LoadConfig() (err error) {
	config, err = storage.LoadConfig()
	return
}

func ConnectToInboxes() (err error) {
	if config == nil {
		return ErrConfigNotLoaded
	}
	for _, p := range config.Providers {
		inbox, err := provider.Factories[p.Type](p.Info)
		if err != nil {
			return err
		}
		inboxes = append(inboxes, inbox)
	}
	return nil
}

func GetInboxes() ([]provider.MailBox, error) {
	if config == nil {
		return nil, ErrConfigNotLoaded
	}
	return inboxes, nil
}
