package server

import (
	"strings"

	"github.com/Pauloo27/mail-notifier/socket/common"
	"github.com/Pauloo27/mail-notifier/socket/common/command"
)

type handlerFunction func(command string, args []string) *common.Response

var commandMap = map[string]handlerFunction{
	command.EchoCommand.Name: echoCommand,
}

func echoCommand(command string, args []string) *common.Response {
	return &common.Response{
		Error: nil,
		Data:  strings.Join(args, " "),
	}
}
