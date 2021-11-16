package config

import "github.com/Pauloo27/mail-notifier/internal/storage"

var (
	Config *storage.Config
)

func Load() (err error) {
	Config, err = storage.LoadConfig()
	return
}
