package main

import (
	"fmt"
	"strconv"

	"github.com/Pauloo27/mail-notifier/cli/polybar"
	"github.com/Pauloo27/mail-notifier/socket/client"
	"github.com/Pauloo27/mail-notifier/socket/common/types"
)

var (
	inboxes []*types.Inbox
)

func printStatus(unreadCount int) {
	coolButton := polybar.ActionButton{
		Index:   polybar.LeftClick,
		Display: " : " + strconv.Itoa(unreadCount),
		Command: "/home/paulo/Dev/Go/src/mail-notifier/mail-notifier-gui",
	}
	fmt.Println(coolButton.String())
}

func handleError(err error) {
	if err != nil {
		panic(err)
	}
}

func mustListInboxes(client *client.Client) {
	var err error
	inboxes, err = client.ListInboxes()
	handleError(err)
}

func mustFetchUnread(client *client.Client) (unread int) {
	for _, inbox := range inboxes {
		unreadMessages, err := client.FetchUnreadMessagesIn(inbox.ID)
		handleError(err)
		unread += len(unreadMessages.Messages)
	}
	return
}

func mustListenToChanges(client *client.Client, ch chan int) {
	for _, inbox := range inboxes {
		err := client.UnlistenToInbox(inbox.ID)
		handleError(err)
	}
}

func main() {
	client := client.NewClient()
	handleError(client.Connect())
	mustListInboxes(client)
	printStatus(mustFetchUnread(client))
	ch := make(chan int)
	mustListenToChanges(client, ch)
	/*
		for {
			printStatus(<-ch)
		}
	*/
}
