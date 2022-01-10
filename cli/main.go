package main

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/Pauloo27/logger"
	"github.com/Pauloo27/mail-notifier/cli/polybar"
	"github.com/Pauloo27/mail-notifier/socket/client"
	"github.com/Pauloo27/mail-notifier/socket/common/types"
)

var (
	inboxes       []*types.Inbox
	unreadByInbox = make(map[int]int)
)

func printStatus(unreadCount int) {
	color := "#50fa7b"
	if unreadCount != 0 {
		color = "#ffb86c"
	}
	coolButton := polybar.ActionButton{
		Index:          polybar.LeftClick,
		Display:        " : " + strconv.Itoa(unreadCount) + " unread",
		UnderlineColor: color,
		Command:        "mail-notifier-gui",
	}
	fmt.Println(coolButton.String())
}

func handleError(err error) {
	if err == nil {
		return
	}
	errBtn := polybar.ActionButton{
		Index:          polybar.LeftClick,
		Display:        "error",
		UnderlineColor: "#ff5555",
	}
	fmt.Println(errBtn)
	logger.Fatal(err)
}

func mustListInboxes(client *client.Client) {
	var err error
	inboxes, err = client.ListInboxes()
	handleError(err)
}

func mustFetchUnread(client *client.Client) (unread int) {
	for i, inbox := range inboxes {
		unreadMessages, err := client.FetchUnreadMessagesIn(inbox.ID)
		handleError(err)
		unreadByInbox[i] = len(unreadMessages.Messages)
		unread += len(unreadMessages.Messages)
	}
	return
}

func mustListenToChanges(c *client.Client, ch chan int) {
	c.OnInboxChanged = func(inboxID int, messages *types.CachedUnreadMessages) {
		unreadByInbox[inboxID] = len(messages.Messages)
		sum := 0
		for _, unread := range unreadByInbox {
			sum += unread
		}
		ch <- sum
	}
	for _, inbox := range inboxes {
		err := c.ListenToInbox(inbox.ID)
		handleError(err)
	}
}

func init() {
	fileName := fmt.Sprintf("/tmp/mail-notifier/log-%d.txt", time.Now().Unix())
	logFile, err := os.Create(fileName)
	if err != nil {
		panic(err)
	}
	logger.Stdout = logFile
	logger.Stderr = logFile
}

func main() {
	client := client.NewClient()
	handleError(client.Connect())
	mustListInboxes(client)
	printStatus(mustFetchUnread(client))
	ch := make(chan int)
	mustListenToChanges(client, ch)
	for {
		printStatus(<-ch)
	}
}
