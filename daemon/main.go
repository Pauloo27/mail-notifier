package main

import (
	"github.com/Pauloo27/logger"
	"github.com/Pauloo27/mail-notifier/socket/server"
	"github.com/Pauloo27/mail-notifier/socket/server/data"

	_ "github.com/Pauloo27/mail-notifier/core/provider/gmail"
	_ "github.com/Pauloo27/mail-notifier/core/provider/mail"
	_ "github.com/emersion/go-message/charset"
)

func main() {
	logger.Info("loading config...")
	if err := data.LoadConfig(); err != nil {
		logger.Fatal(err)
	}
	logger.Success("config loaded!")
	logger.Info("connecting to inboxes...")
	if err := data.ConnectToInboxes(); err != nil {
		logger.Fatal(err)
	}
	logger.Success("connected to all inboxes!")
	sv := server.NewServer()
	if err := sv.Listen(); err != nil {
		logger.Fatal(err)
	}
}
