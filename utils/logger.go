package utils

import (
	"fmt"
	"io/ioutil"
	"log"
	"strings"
	"time"
)

const StatusFile = "/dev/shm/gmail-status.txt"

func HandleFatal(message string, err error) {
	if err != nil {
		fmt.Println(message)
		LogErrorStatus(message)
		log.Fatalf("[FATAL]: %s: %v", message, err)
	}
}

func LogErrorStatus(message string) {
	now := time.Now()
	timestamp := fmt.Sprintf("Checked at %s", now.Format("Mon, 2 Jan • 15:04"))
	unreadCount := "-"
	data := []byte(fmt.Sprintf("%s\n%s\n%s\n", unreadCount, timestamp, message))
	err := ioutil.WriteFile(StatusFile, data, 0644)
	if err != nil {
		HandleFatal("Cannot write file "+StatusFile, err)
	}
}

func LogStatus(status []int) {
	now := time.Now()
	timestamp := fmt.Sprintf("Checked at %s", now.Format("Mon, 2 Jan • 15:04"))
	unreadCount := ""
	for _, unread := range status {
		unreadCount += fmt.Sprintf("%d ", unread)
	}
	data := []byte(fmt.Sprintf("%s\n%s\n", strings.TrimSuffix(unreadCount, " "), timestamp))
	err := ioutil.WriteFile(StatusFile, data, 0644)
	if err != nil {
		HandleFatal("Cannot write file "+StatusFile, err)
	}
}
