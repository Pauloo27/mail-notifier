package server

import (
	"strings"

	"github.com/Pauloo27/mail-notifier/socket/common"
	"github.com/Pauloo27/mail-notifier/socket/common/command"
	"github.com/Pauloo27/mail-notifier/socket/server/data"
)

type handlerFunction func(command string, args []string) *common.Response

var commandMap = map[string]handlerFunction{
	command.EchoCommand.Name:        echoCommand,
	command.ListInboxesCommand.Name: listInboxes,
}

func echoCommand(command string, args []string) *common.Response {
	return &common.Response{
		Error: nil,
		Data:  strings.Join(args, " "),
	}
}

func listInboxes(command string, args []string) *common.Response {
	inboxes, err := data.GetInboxes()

	var d []map[string]interface{}

	for _, inbox := range inboxes {
		d = append(d, map[string]interface{}{
			"address": inbox.GetAddress(),
		})
	}

	return &common.Response{
		Error: err,
		Data:  d,
	}
}
