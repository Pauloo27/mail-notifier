package client

import (
	"encoding/json"
	"errors"
	"fmt"
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

type Client struct {
	OnInboxChanged func(inboxID int, messages *types.CachedUnreadMessages)
	LastInboxList  []*types.Inbox
}

func NewClient() *Client {
	return &Client{}
}

func (c *Client) handler(req *common.Request) (interface{}, error) {
	if req.Command == command.NotifyInboxChange.Name {
		if len(req.Args) != 1 {
			return nil, fmt.Errorf("invalid argument size: %d", len(req.Args))
		}

		inboxID, err := strconv.Atoi(req.Args[0])
		if err != nil {
			return nil, fmt.Errorf("invalid inbox id: %w", err)
		}

		if c.OnInboxChanged == nil {
			return "ignored, but ok", nil
		}

		var unread types.CachedUnreadMessages
		rawData, err := json.Marshal(req.Data)
		if err != nil {
			return nil, err
		}
		err = json.Unmarshal(rawData, &unread)
		if err != nil {
			return nil, err
		}

		c.OnInboxChanged(inboxID, &unread)
		return "ok", nil
	}
	return nil, errors.New("command not handled")
}

func (c *Client) Connect() error {
	var err error
	conn, err = net.Dial("unix", common.SocketPath)
	if err != nil {
		return err
	}
	t = transport.NewTransport(conn)
	go t.Start(c.handler)
	go t.TransmitHeartbeats()
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
	if err == nil {
		c.LastInboxList = inboxes
	}
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

func (c *Client) ClearAllInboxCache() error {
	res, err := c.sendCommand(command.ClearAllInboxesCache.Name, nil)
	if err != nil {
		return err
	}
	return res.Error
}

func (c *Client) ListenToInbox(inboxID int) error {
	res, err := c.sendCommand(command.ListenToInbox.Name, []string{strconv.Itoa(inboxID)})
	if err != nil {
		return err
	}
	return res.Error
}

func (c *Client) UnlistenToInbox(inboxID int) error {
	res, err := c.sendCommand(command.UnlistenToInbox.Name, []string{strconv.Itoa(inboxID)})
	if err != nil {
		return err
	}
	return res.Error
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
	t.Send(command, args, nil, cb)
	res := <-resCh
	return res, nil
}
