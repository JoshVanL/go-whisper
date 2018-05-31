package server

import (
	"crypto/rsa"
	"fmt"
	"net"

	"github.com/sirupsen/logrus"

	"github.com/joshvanl/go-whisper/pkg/key"
)

type Server struct {
	log *logrus.Entry

	addr string
	dir  string

	key  *rsa.PrivateKey
	conn net.Conn

	uids map[uint64]*rsa.PublicKey
}

func New(addr string, dir string, log *logrus.Entry) (*Server, error) {

	k, err := key.RetrieveLocalKey(dir)
	if err != nil {
		return nil, fmt.Errorf("failed to read server key: %v", err)
	}

	pubKeys, err := key.RetrieveUIDPublicKeys(dir)
	if err != nil {
		return nil, fmt.Errorf("failed to read local client keys: %v", err)
	}

	return &Server{
			log:  log,
			addr: addr,
			dir:  dir,
			key:  k,
			uids: pubKeys,
		},
		nil

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
