package server

import (
	"bufio"
	"errors"
	"net"
	"os"

	"github.com/Pauloo27/logger"
	"github.com/Pauloo27/mail-notifier/socket/common"
	"github.com/Pauloo27/mail-notifier/socket/common/transport"
)

type Server struct {
}

func NewServer() *Server {
	return &Server{}
}

func handleCommand(command string, args []string) (interface{}, error) {
	handler, ok := commandMap[command]
	if !ok {
		return nil, errors.New("command not found")
	}
	return handler(command, args)
}

func (s *Server) handleConnection(conn net.Conn) error {
	rw := bufio.NewReadWriter(bufio.NewReader(conn), bufio.NewWriter(conn))
	transport := transport.NewTransport(rw)
	return transport.Start(func(req *common.Request) (interface{}, error) {
		return handleCommand(req.Command, req.Args)
	})
}

func (s *Server) acceptNewConnections(l net.Listener) error {
	for {
		conn, err := l.Accept()
		if err != nil {
			return err
		}
		go func() {
			err := s.handleConnection(conn)
			if err != nil {
				logger.Error(err)
			}
		}()
	}
}

func (s *Server) Listen() error {
	os.MkdirAll(common.SocketPathRootDir, 0700)
	if _, err := os.Stat(common.SocketPath); !os.IsNotExist(err) {
		os.Remove(common.SocketPath)
	}
	l, err := net.Listen("unix", common.SocketPath)
	if err != nil {
		return err
	}
	return s.acceptNewConnections(l)
}
