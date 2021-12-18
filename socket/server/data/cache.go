package data

import (
	"crypto/md5"
	"sort"
	"strings"
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

func refreshCache(inboxID int) {
	prevSum := calcCurrentInboxMD5(inboxID)
	for {
		curSum := calcCurrentInboxMD5(inboxID)
		<-cacheTimers[inboxID].C
		ClearInboxCache(inboxID)
		err := fetchUnreadMessage(inboxID)
		if err != nil {
			panic(err)
		}
		if curSum != prevSum {
			// TODO: notify listeners
		}
		prevSum = curSum
	}
}
