package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"time"

	"github.com/Pauloo27/gmail-notifier/gapi"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/gmail/v1"
)

const statusFile = "/dev/shm/gmail-status.txt"

func fetchMessages(srv *gmail.Service) []*gmail.Message {
	r, err := srv.Users.Messages.List("me").LabelIds("UNREAD").IncludeSpamTrash(true).MaxResults(10).Do()
	if err != nil {
		log.Fatalf("Unable to retrieve messages: %v", err)
	}

	return r.Messages
}

func logStatus(unreadCount int) {
	now := time.Now()
	timestamp := fmt.Sprintf("Checked at %s", now.Format("Mon, 2 Jan • 15:04"))
	data := []byte(fmt.Sprintf("%d\n%s\n", unreadCount, timestamp))
	err := ioutil.WriteFile(statusFile, data, 0644)
	if err != nil {
		log.Fatalf("Cannot write file %s: %v", statusFile, err)
	}
}

func main() {
	b, err := ioutil.ReadFile("secret/credentials.json")
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}

	config, err := google.ConfigFromJSON(b, gmail.GmailReadonlyScope)
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}
	client := gapi.GetClient(config)

	srv, err := gmail.New(client)
	if err != nil {
		log.Fatalf("Unable to retrieve Gmail client: %v", err)
	}

	fmt.Printf("Started. Logging status to %s\n", statusFile)

	for {
		messages := fetchMessages(srv)
		messageCount := len(messages)
		logStatus(messageCount)

		time.Sleep(1 * time.Minute)
	}
}
