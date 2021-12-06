package client

import (
	"bufio"
	"encoding/json"
	"net"

	"github.com/Pauloo27/mail-notifier/socket/common"
	"github.com/Pauloo27/mail-notifier/socket/common/command"
	"github.com/Pauloo27/mail-notifier/socket/common/transport"
	"github.com/Pauloo27/mail-notifier/socket/common/types"
)

var (
	conn net.Conn
	t    *transport.Transport
)

func Connect() error {
	var err error
	conn, err = net.Dial("unix", common.SocketPath)
	if err != nil {
		return err
	}
	rw := bufio.NewReadWriter(bufio.NewReader(conn), bufio.NewWriter(conn))
	t = transport.NewTransport(rw)
	go t.Start(nil)
	return nil
}

func ListInboxes() ([]*types.Inbox, error) {
	res, err := SendCommand(command.ListInboxesCommand.Name, nil)
	if err != nil {
		return nil, err
	}
	var inboxes []*types.Inbox
	rawData, err := json.Marshal(res.Data)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(rawData, &inboxes)
	return inboxes, err
}

func SendCommand(command string, args []string) (*common.Response, error) {
	resCh := make(chan *common.Response)
	cb := func(res *common.Response) {
		resCh <- res
	}
	t.Send(command, args, cb)
	res := <-resCh
	return res, nil
}
