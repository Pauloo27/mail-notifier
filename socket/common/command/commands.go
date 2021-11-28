package command

var EchoCommand = &Command{
	Name:        "echo",
	Usage:       "echo <data>",
	Description: "send back the received data",
}

var ListInboxesCommand = &Command{
	Name:        "list_inboxes",
	Usage:       "list_inboxes",
	Description: "list the inboxes",
}
