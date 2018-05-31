package server

import (
	"crypto/rsa"
	"fmt"
	"net"

	"github.com/sirupsen/logrus"

	"github.com/joshvanl/go-whisper/pkg/config"
	"github.com/joshvanl/go-whisper/pkg/key"
)

type Server struct {
	log *logrus.Entry

	addr string
	dir  string

	clientUids map[string]bool
	uids       *key.UIDs

	key  *rsa.PrivateKey
	conn net.Conn

	config *config.Config
}

func New(addr string, dir string, log *logrus.Entry) (*Server, error) {

	log.Infof("Retrieving local key pair...")
	k, err := key.RetrieveLocalKey(dir)
	if err != nil {
		return nil, fmt.Errorf("failed to read server key: %v", err)
	}

	server := &Server{
		log:  log,
		addr: addr,
		dir:  dir,
		key:  k,
	}

	log.Infof("Retrieving local server config...")
	config, err := config.ReadConfig(dir)
	if err != nil {
		return nil, fmt.Errorf("failed to read config: %v", err)
	}
	server.config = config

	if addr != "" {
		server.addr = addr
	} else {
		server.addr = server.config.Address
	}

	log.Infof("Retrieving local uids...")
	uids, err := key.NewUIDs(server.dir, 0)
	if err != nil {
		return nil, err
	}
	server.uids = uids

	return server, nil
}

func (s *Server) Serve() error {
	ln, err := net.Listen("tcp", s.addr)
	if err != nil {
		return fmt.Errorf("failed to serve address: %v", err)
	}

	for {

		con, err := ln.Accept()
		if err != nil {
			fmt.Printf("failed to accept connection: %v", err)
			continue
		}

		go s.Handle(con)
	}

}
