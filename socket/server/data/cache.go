package data

func refreshCache(inboxID int) {
	for {
		<-cacheTimers[inboxID].C
		ClearInboxCache(inboxID)
		fetchUnreadMessage(inboxID)
	}
}
