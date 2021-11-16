package main

import (
	"github.com/Pauloo27/mail-notifier/gui/internal/config"
	"github.com/Pauloo27/mail-notifier/gui/internal/containers/home"

	_ "github.com/Pauloo27/mail-notifier/internal/providers/mail"
)

func main() {
	err := config.Load()
	if err != nil {
		// TODO: create file when not found...
		panic(err)
	}

	home.Show()
}
