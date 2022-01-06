package server

import (
	"errors"
	"io"
	"net"
	"os"
	"strconv"

	"github.com/Pauloo27/logger"
	"github.com/Pauloo27/mail-notifier/socket/common"
	"github.com/Pauloo27/mail-notifier/socket/common/command"
	"github.com/Pauloo27/mail-notifier/socket/common/transport"
	"github.com/Pauloo27/mail-notifier/socket/common/types"
	"github.com/Pauloo27/mail-notifier/socket/server/data"
)

type Server struct {
	clients map[string]*ConnectedClient
}

type ConnectedClient struct {
	Transport *transport.Transport
}

func NewServer() *Server {
	server := &Server{
		clients: make(map[string]*ConnectedClient),
	}
	data.NotifyInboxChanges = func(clientID string, inboxID int, messages *types.CachedUnreadMessages) {
		client, ok := server.clients[clientID]
		if !ok {
			return
		}
		client.Transport.Send(command.NotifyInboxChange.Name, []string{strconv.Itoa(inboxID)}, messages, nil)
	}
	return server
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
	s.clients[transport.UID] = client
	logger.Infof("new client connected: %s", transport.UID)
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
			client, err := s.handleConnection(conn)
			logger.Infof("client disconnected: %s", client.Transport.UID)
			delete(s.clients, client.Transport.UID)
			_ = client.Transport.Stop()
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
