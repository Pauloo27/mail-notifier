package server

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/Pauloo27/mail-notifier/socket/common"
	"github.com/Pauloo27/mail-notifier/socket/common/command"
	"github.com/Pauloo27/mail-notifier/socket/server/data"
)

type handlerFunction func(command string, args []string) *common.Response

var commandMap = map[string]handlerFunction{
	command.EchoCommand.Name:         echoCommand,
	command.ListInboxesCommand.Name:  listInboxes,
	command.FetchMessageCommand.Name: fetchMessage,
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

func fetchMessage(command string, args []string) *common.Response {
	if len(args) != 2 {
		return &common.Response{
			Error: fmt.Errorf("invalid argument size: %d", len(args)),
			Data:  nil,
		}
	}
	inboxID, err := strconv.Atoi(args[0])
	if err != nil {
		return &common.Response{
			Error: fmt.Errorf("invalid inbox id: %w", err),
			Data:  nil,
		}
	}

	msg, err := data.GetMessage(inboxID, args[1])
	if err != nil {
		return &common.Response{
			Error: err,
			Data:  nil,
		}
	}

	d := map[string]interface{}{
		"id":      (*msg).GetID(),
		"to":      (*msg).GetTo(),
		"from":    (*msg).GetFrom(),
		"date":    (*msg).GetDate(),
		"subject": (*msg).GetSubject(),
	}

	return &common.Response{
		Error: err,
		Data:  d,
	}
}
