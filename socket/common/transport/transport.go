package transport

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"sync"
	"time"

	"github.com/Pauloo27/logger"
	"github.com/Pauloo27/mail-notifier/socket/common"
	"github.com/google/uuid"
)

type Transport struct {
	UID             string
	conn            net.Conn
	health          *Health
	rw              *bufio.ReadWriter
	writeLock       *sync.Mutex
	pendingRequests map[string]ResponseCallback
}

func (t *Transport) Stop() error {
	t.health.Kill()
	return t.conn.Close()
}

func NewTransport(conn net.Conn) *Transport {
	rw := bufio.NewReadWriter(bufio.NewReader(conn), bufio.NewWriter(conn))
	t := &Transport{
		UID:             uuid.NewString(),
		conn:            conn,
		rw:              rw,
		writeLock:       &sync.Mutex{},
		pendingRequests: make(map[string]ResponseCallback),
	}
	t.health = newHealth(func() { _ = t.Stop() })
	return t
}

func (t *Transport) TransmitHeartbeats() {
	heartbeat := "<3"
	for !t.health.dead {
		time.Sleep(heartbeatRate)
		_, err := t.Send(heartbeatCommandName, []string{heartbeat}, nil, func(res *common.Response) {
			heartbeatRes, ok := res.Data.(string)
			if !ok || heartbeatRes != heartbeat {
				logger.Fatal("invalid heartbeat response")
				t.Stop()
				return
			}
			t.health.HeartbeatReceived()
		})
		if err != nil {
			logger.Fatal("cannot send heartbeat: ", err)
			t.Stop()
			break
		}
		t.health.HeartbeatSent()
	}
}

type RequestHandler func(req *common.Request) (data interface{}, err error)
type ResponseCallback func(res *common.Response)

func (t *Transport) Start(handler RequestHandler) error {
	if handler == nil {
		return nil
	}
	return t.doRead(handler)
}

func (t *Transport) writeFullPackage(data []byte) error {
	t.writeLock.Lock()
	defer t.writeLock.Unlock()
	w, err := t.rw.Write(data)
	if err != nil {
		return err
	}
	if w != len(data) {
		return errors.New("wrong data written")
	}
	err = t.rw.WriteByte('\n')
	if err != nil {
		return err
	}
	return t.rw.Flush()
}

func (t *Transport) Send(command string, args []string, data interface{}, cb ResponseCallback) (string, error) {
	id := uuid.New().String()
	req := &common.Request{
		ID:      id,
		Command: command,
		Data:    data,
		Args:    args,
	}
	reqJSON, err := json.Marshal(req)
	if err != nil {
		return "", err
	}
	t.pendingRequests[id] = cb
	err = t.writeFullPackage(reqJSON)
	if err != nil {
		cb(&common.Response{
			Error: fmt.Errorf("failed to send: %v", err),
			To:    id,
			Data:  nil,
		})
		return "", err
	}
	return id, nil
}

func (t *Transport) Respond(req *common.Request, data interface{}, err error) error {
	res := &common.Response{
		Data:  data,
		Error: err,
		To:    req.ID,
	}
	resJSON, err := json.Marshal(res)
	if err != nil {
		return err
	}
	return t.writeFullPackage(resJSON)
}

func isReq(req common.Request) bool {
	return req.ID != "" && req.Command != ""
}

func isRes(res common.Response) bool {
	return res.To != ""
}

func (t *Transport) doRead(handler RequestHandler) error {
	for {
		data, err := t.rw.ReadString('\n')
		if err != nil {
			return err
		}
		var req common.Request
		err = json.Unmarshal([]byte(data), &req)
		if err == nil && isReq(req) {
			if req.Command == heartbeatCommandName {
				t.health.HeartbeatReceived()
				go t.Respond(&req, req.Args[0], nil)
				continue
			}
			go func() {
				data, err := handler(&req)
				t.Respond(&req, data, err)
			}()
			continue
		}
		var res common.Response
		err = json.Unmarshal([]byte(data), &res)
		if err == nil && isRes(res) {
			go func() {
				cb, found := t.pendingRequests[res.To]
				if !found || cb == nil {
					return
				}
				cb(&res)
			}()
			continue

		}
		logger.Warn("invalid data received:", data)
	}
}
