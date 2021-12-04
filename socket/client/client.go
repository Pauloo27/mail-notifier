package client

import (
	"bufio"
	"encoding/json"
	"errors"
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
	rw := bufio.NewReadWriter(bufio.NewReader(conn), bufio.NewWriter(conn))
	t = transport.NewTransport(rw)
	return err
}

func ListInboxes() ([]*types.Inbox, error) {
	res, err := SendCommand(command.ListInboxesCommand.Name, nil)
	if err != nil {
		return nil, err
	}
	dataList, ok := res.Data.([]map[string]interface{})
	if !ok {
		return nil, errors.New("unexpected response")
	}
	var inboxes []*types.Inbox
	for _, data := range dataList {
		jsonData, err := json.Marshal(data)
		if err != nil {
			return nil, err
		}
		var inbox types.Inbox
		err = json.Unmarshal(jsonData, &inbox)
		if err != nil {
			return nil, err
		}
		inboxes = append(inboxes, &inbox)
	}
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
