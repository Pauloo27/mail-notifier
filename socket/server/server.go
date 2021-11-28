package server

import (
	"bufio"
	"encoding/json"
	"errors"
	"net"
	"os"
	"strings"

	"github.com/Pauloo27/mail-notifier/socket/common"
)

func handleCommand(line string) *common.Response {
	line = strings.ToLower(line)
	parts := strings.Split(line, " ")
	command := parts[0]
	var args []string
	if len(parts) > 1 {
		args = parts[1:]
	}
	handler, ok := commandMap[command]
	if !ok {
		return &common.Response{
			Error: errors.New("command not found"),
		}
	}
	return handler(command, args)
}

func handleConnection(conn net.Conn) error {
	rw := bufio.NewReadWriter(bufio.NewReader(conn), bufio.NewWriter(conn))
	for {
		line, err := rw.ReadString('\n')
		if err != nil {
			return err
		}
		response := handleCommand(strings.TrimSuffix(line, "\n"))
		jsonResponse, err := json.Marshal(response)
		if err != nil {
			return err
		}
		jsonResponse = append(jsonResponse, '\n')
		if _, err = rw.Write(jsonResponse); err != nil {
			return err
		}
		if err = rw.Flush(); err != nil {
			return err
		}
	}
}

func acceptNewConnections(l net.Listener) error {
	for {
		conn, err := l.Accept()
		if err != nil {
			return err
		}
		go handleConnection(conn) // TODO: handle error?
	}
}

func Listen() error {
	os.MkdirAll(common.SocketPathRootDir, 0700)
	if _, err := os.Stat(common.SocketPath); !os.IsNotExist(err) {
		os.Remove(common.SocketPath)
	}
	l, err := net.Listen("unix", common.SocketPath)
	if err != nil {
		return err
	}
	return acceptNewConnections(l)
}
