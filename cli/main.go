package main

import (
	"fmt"
	"strconv"
	"time"

	"github.com/Pauloo27/mail-notifier/cli/polybar"
	"github.com/Pauloo27/mail-notifier/socket/client"
)

func printStatus(unreadCount int) {
	coolButton := polybar.ActionButton{
		Index:   polybar.LEFT_CLICK,
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

func mustFetchUnread(client *client.Client) (unread int) {
	inboxes, err := client.ListInboxes()
	handleError(err)
	for _, inbox := range inboxes {
		unreadMessages, err := client.FetchUnreadMessagesIn(inbox.ID)
		handleError(err)
		unread += len(unreadMessages.Messages)
	}
	return
}

func mustListenToChanges(client *client.Client, ch chan int) {
	go func() {
		time.Sleep(10 * time.Second)
		ch <- 0
	}()
}

func main() {
	client := client.NewClient()
	handleError(client.Connect())
	printStatus(mustFetchUnread(client))
	ch := make(chan int)
	mustListenToChanges(client, ch)
	for {
		printStatus(<-ch)
	}
}
