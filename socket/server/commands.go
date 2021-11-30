package server

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/Pauloo27/mail-notifier/socket/common/command"
	"github.com/Pauloo27/mail-notifier/socket/common/types"
	"github.com/Pauloo27/mail-notifier/socket/server/data"
)

type handlerFunction func(command string, args []string) (data interface{}, err error)

var commandMap = map[string]handlerFunction{
	command.EchoCommand.Name:         echoCommand,
	command.ListInboxesCommand.Name:  listInboxes,
	command.FetchMessageCommand.Name: fetchMessage,
}

func echoCommand(command string, args []string) (interface{}, error) {
	return strings.Join(args, " "), nil
}

func listInboxes(command string, args []string) (interface{}, error) {
	inboxes, err := data.GetInboxes()

	var i types.Inboxes

	for _, inbox := range inboxes {
		i = append(i, &types.Inbox{
			MailBox: inbox,
		})
	}

	return i, err
}

func fetchMessage(command string, args []string) (interface{}, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf("invalid argument size: %d", len(args))
	}

	inboxID, err := strconv.Atoi(args[0])
	if err != nil {
		return nil, fmt.Errorf("invalid inbox id: %w", err)
	}

	msg, err := data.GetMessage(inboxID, args[1])

	if err != nil {
		return nil, err
	}

	return msg, err
}
