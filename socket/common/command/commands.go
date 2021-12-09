package command

var EchoCommand = &Command{
	Name:        "echo",
	Usage:       "echo [data]",
	Description: "send back the received data",
}

var ListInboxesCommand = &Command{
	Name:        "list_inboxes",
	Usage:       "list_inboxes",
	Description: "list the inboxes",
}

var FetchMessageCommand = &Command{
	Name:        "fetch_message",
	Usage:       "fetch_message [inbox id] [message id]",
	Description: "fetch a message from an inbox",
}

var FetchUnreadMessagesIn = &Command{
	Name:        "fetch_unread_messages_in",
	Usage:       "fetch_unread_messages_in [inbox id]",
	Description: "fetch unread messages from an inbox",
}

var FetchAllUnreadMessages = &Command{
	Name:        "fetch_all_unread_messages",
	Usage:       "fetch_all_unread_messages [inbox id]",
	Description: "fetch unread messages from all inboxes",
}

var MarkMessageAsRead = &Command{
	Name:        "mark_message_as_read",
	Usage:       "mark_message_as_read [inbox id] [message id]",
	Description: "mark a message as read",
}

var ClearInboxCache = &Command{
	Name:        "clear_inbox_cache",
	Usage:       "clear_inbox_cache [inbox id]",
	Description: "clear an inbox unread messages cache",
}

var ClearAllInboxesCache = &Command{
	Name:        "clear_all_inboxes_cache",
	Usage:       "clear_all_inboxes_cache",
	Description: "clear all unread messages cache",
}
