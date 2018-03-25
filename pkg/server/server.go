package server

import (
	"fmt"
	"crypto/rsa"
	"net"

	"github.com/sirupsen/logrus"

	"github.com/joshvanl/go-whisper/pkg/file"
)

type Server struct {
	log  *logrus.Entry
	addr string

	key  *rsa.PrivateKey
	conn net.Conn
}

func New(addr string, log *logrus.Entry) (*Server, error) {

	k, err := file.RetrieveKey()
	if err != nil {
		return nil, fmt.Errorf("failed to read server key: %v", err)
	}

	return &Server{
		log: log,
		addr: addr,
		key: k,
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

		go s.handle(con)
	}

}

func (s *Server) handle(con net.Conn) {
	buff := make([]byte, 256)

	n, err := con.Read(buff)
	if err != nil {
		return
	}

	fmt.Printf("(%d) read from connection: %s\n", n, string(buff))
}
