package transport

import (
	"bufio"
	"encoding/json"
	"errors"
	"sync"

	"github.com/Pauloo27/logger"
	"github.com/Pauloo27/mail-notifier/socket/common"
	"github.com/google/uuid"
)

type Transport struct {
	rw              *bufio.ReadWriter
	writeLock       *sync.Mutex
	pendingRequests map[string]ResponseCallback
}

func NewTransport(rw *bufio.ReadWriter) *Transport {
	return &Transport{
		rw:        rw,
		writeLock: &sync.Mutex{},
	}
}

type RequestHandler func(req *common.Request) (data interface{}, err error)
type ResponseCallback func(res *common.Response)

func (t *Transport) Start(handler RequestHandler) error {
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

func (t *Transport) Send(command string, args []string, cb ResponseCallback) (string, error) {
	id := uuid.New().String()
	req := &common.Request{
		ID:      id,
		Command: command,
		Args:    args,
	}
	reqJSON, err := json.Marshal(req)
	if err != nil {
		return "", err
	}
	err = t.writeFullPackage(reqJSON)
	if err != nil {
		return "", err
	}
	t.pendingRequests[id] = cb
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
				if !found {
					return
				}
				cb(&res)
			}()
			continue

		}
		logger.Warn("invalid data received")
		t.writeFullPackage([]byte("?"))
	}
}
