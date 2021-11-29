package main

import (
	"fmt"

	"github.com/Pauloo27/mail-notifier/socket/server"
	"github.com/Pauloo27/mail-notifier/socket/server/data"

	_ "github.com/Pauloo27/mail-notifier/core/provider/gmail"
	_ "github.com/Pauloo27/mail-notifier/core/provider/mail"
)

func main() {
	if err := data.LoadConfig(); err != nil {
		panic(err)
	}
	if err := data.ConnectToInboxes(); err != nil {
		panic(err)
	}
	fmt.Println("Config loaded!")
	if err := server.Listen(); err != nil {
		panic(err)
	}
}
