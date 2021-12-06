package main

import (
	"github.com/Pauloo27/mail-notifier/gui/internal/containers/home"
	"github.com/Pauloo27/mail-notifier/socket/client"

	_ "github.com/Pauloo27/mail-notifier/core/provider/gmail"
	_ "github.com/Pauloo27/mail-notifier/core/provider/mail"
)

func main() {
	c := client.NewClient()
	err := c.Connect()
	if err != nil {
		panic(err)
	}
	home.Show(c)
}
