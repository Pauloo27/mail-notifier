package server

import (
	"errors"
	"io"
	"net"
	"os"

	"github.com/Pauloo27/logger"
	"github.com/Pauloo27/mail-notifier/socket/common"
	"github.com/Pauloo27/mail-notifier/socket/common/transport"
)

type Server struct {
	clients []*ConnectedClient
}

type ConnectedClient struct {
	Transport *transport.Transport
}

func NewServer() *Server {
	return &Server{}
}

func handleCommand(c *ConnectedClient, command string, args []string) (interface{}, error) {
	handler, ok := commandMap[command]
	if !ok {
		return nil, errors.New("command not found")
	}
	return handler(c, command, args)
}

func (s *Server) handleConnection(conn net.Conn) (*ConnectedClient, error) {
	transport := transport.NewTransport(conn)
	client := &ConnectedClient{transport}
	s.clients = append(s.clients, client)
	return client, transport.Start(func(req *common.Request) (interface{}, error) {
		return handleCommand(client, req.Command, req.Args)
	})
}

func (s *Server) acceptNewConnections(l net.Listener) error {
	for {
		conn, err := l.Accept()
		if err != nil {
			return err
		}
		go func() {
			_, err := s.handleConnection(conn)
			_ = conn.Close()
			if err != nil && !errors.Is(err, io.EOF) {
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
