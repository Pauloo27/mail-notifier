package data

import (
	"crypto/md5"
	"sort"
	"strings"
	"time"

	"github.com/Pauloo27/logger"
)

const (
	refreshCacheAfter = 5 * time.Minute
)

func calcCurrentInboxMD5(inboxID int) [16]byte {
	var msgsIDs []string
	for _, msg := range unreadMessages[inboxID].Messages {
		msgsIDs = append(msgsIDs, msg.ID)
	}
	sort.Strings(msgsIDs)
	x := []byte(strings.Join(msgsIDs, " "))
	return md5.Sum(x)
}

func notifyListeners(inboxID int) {
	if NotifyInboxChanges == nil {
		return
	}
	listeners, found := inboxListeners[inboxID]
	if !found {
		return
	}
	messages, found := unreadMessages[inboxID]
	if !found {
		return
	}
	for _, clientID := range listeners {
		NotifyInboxChanges(clientID, inboxID, messages)
	}
}

func refreshCache(inboxID int) {
	prevSum := calcCurrentInboxMD5(inboxID)
	for {
		<-cacheTimers[inboxID].C
		ClearInboxCache(inboxID)
		err := fetchUnreadMessage(inboxID)
		curSum := calcCurrentInboxMD5(inboxID)
		if err != nil {
			logger.Fatal(err)
		}
		if curSum != prevSum {
			notifyListeners(inboxID)
		}
		prevSum = curSum
	}
}
