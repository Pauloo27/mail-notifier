package main

import (
	"github.com/Pauloo27/mail-notifier/socket/server"
)

func main() {
	if err := server.Listen(); err != nil {
		panic(err)
	}
}
