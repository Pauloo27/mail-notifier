package client

import (
	"encoding/json"
	"net"
	"strconv"

	"github.com/Pauloo27/mail-notifier/socket/common"
	"github.com/Pauloo27/mail-notifier/socket/common/command"
	"github.com/Pauloo27/mail-notifier/socket/common/transport"
	"github.com/Pauloo27/mail-notifier/socket/common/types"
)

var (
	conn net.Conn
	t    *transport.Transport
)

type Client struct{}

func NewClient() *Client {
	return &Client{}
}

func (c *Client) Connect() error {
	var err error
	conn, err = net.Dial("unix", common.SocketPath)
	if err != nil {
		return err
	}
	t = transport.NewTransport(conn)
	go t.Start(nil)
	return nil
}

func (c *Client) ListInboxes() ([]*types.Inbox, error) {
	res, err := c.sendCommand(command.ListInboxesCommand.Name, nil)
	if err != nil {
		return nil, err
	}
	if res.Error != nil {
		return nil, res.Error
	}
	var inboxes []*types.Inbox
	rawData, err := json.Marshal(res.Data)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(rawData, &inboxes)
	return inboxes, err
}

func (c *Client) FetchUnreadMessagesIn(inboxID int) (*types.CachedUnreadMessages, error) {
	res, err := c.sendCommand(command.FetchUnreadMessagesIn.Name, []string{strconv.Itoa(inboxID)})
	if err != nil {
		return nil, err
	}
	if res.Error != nil {
		return nil, res.Error
	}
	var unread types.CachedUnreadMessages
	rawData, err := json.Marshal(res.Data)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(rawData, &unread)
	return &unread, err
}

func (c *Client) MarkMessageAsRead(inboxID int, messageID string) error {
	res, err := c.sendCommand(command.MarkMessageAsRead.Name, []string{strconv.Itoa(inboxID), messageID})
	if err != nil {
		return err
	}
	return res.Error
}

func (c *Client) sendCommand(command string, args []string) (*common.Response, error) {
	resCh := make(chan *common.Response)
	cb := func(res *common.Response) {
		resCh <- res
	}
	t.Send(command, args, cb)
	res := <-resCh
	return res, nil
}
