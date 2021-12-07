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
	command.EchoCommand.Name:            echoCommand,
	command.ListInboxesCommand.Name:     listInboxes,
	command.FetchMessageCommand.Name:    fetchMessage,
	command.FetchUnreadMessagesIn.Name:  fetchUnreadMessagesIn,
	command.FetchAllUnreadMessages.Name: fetchAllUnreadMessages,
	command.MarkMessageAsRead.Name:      markMessageAsRead,
}

func echoCommand(command string, args []string) (interface{}, error) {
	return strings.Join(args, " "), nil
}

func listInboxes(command string, args []string) (interface{}, error) {
	inboxes, err := data.GetInboxes()

	var i types.Inboxes

	for _, inbox := range inboxes {
		i = append(i, &types.Inbox{
			Address: inbox.GetAddress(),
			WebURL:  inbox.GetWebURL(),
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

func fetchAllUnreadMessages(command string, args []string) (interface{}, error) {
	msgs, err := data.GetAllUnreadMessages()
	if err != nil {
		return nil, err
	}
	return msgs, err
}

func markMessagesAsRead(command string, args []string) (interface{}, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf("invalid argument size: %d", len(args))
	}

	inboxID, err := strconv.Atoi(args[0])
	if err != nil {
		return nil, fmt.Errorf("invalid inbox id: %w", err)
	}

	return nil, data.MarkMessageAsRead(inboxID, args[1])
}

func fetchUnreadMessagesIn(command string, args []string) (interface{}, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("invalid argument size: %d", len(args))
	}

	inboxID, err := strconv.Atoi(args[0])
	if err != nil {
		return nil, fmt.Errorf("invalid inbox id: %w", err)
	}

	msgs, err := data.GetUnreadMessagesIn(inboxID)

	if err != nil {
		return nil, err
	}

	return msgs, err
}
