package client

import (
	"bufio"
	"encoding/json"
	"errors"
	"net"
	"sync"

	"github.com/Pauloo27/mail-notifier/socket/common"
	"github.com/Pauloo27/mail-notifier/socket/common/command"
	"github.com/Pauloo27/mail-notifier/socket/common/types"
)

var (
	conn net.Conn
	rw   *bufio.ReadWriter
	lock *sync.Mutex
)

func Connect() error {
	var err error
	conn, err = net.Dial("unix", common.SocketPath)
	rw = bufio.NewReadWriter(bufio.NewReader(conn), bufio.NewWriter(conn))
	lock = &sync.Mutex{}
	return err
}

func ListInboxes() ([]*types.Inbox, error) {
	res, err := SendCommand(command.ListInboxesCommand.Name)
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

func SendCommand(command string) (*common.Response, error) {
	lock.Lock()
	defer lock.Unlock()
	data := []byte(command)
	w, err := rw.Write(data)
	if err != nil {
		return nil, err
	}
	if w != len(data) {
		return nil, errors.New("something went wrong while writing")
	}
	responseData, err := rw.ReadString('\n')
	if err != nil {
		return nil, err
	}
	var response common.Response
	err = json.Unmarshal([]byte(responseData), &response)
	return &response, err
}
