package common

import "fmt"

var (
	SocketPathRootDir = "/tmp/mail-notifier"
	SocketPath        = fmt.Sprintf("%s/data.sock", SocketPathRootDir)
)
