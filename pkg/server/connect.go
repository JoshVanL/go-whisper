package server

import (
	"crypto/rand"
	"fmt"
	"math"
	"math/big"
	"net"
)

func (s *Server) Handle(con net.Conn) {
	buff := make([]byte, 256)

	n, err := con.Read(buff)
	if err != nil {
		return
	}

	uid, err := s.newUID()
	if err != nil {
		fmt.Errorf("failed to create new uid: %v", err)
	}

	fmt.Printf("buff: %s\nuid: %v\n", buff[0:n], uid)

	//fmt.Printf("(%d) read from connection: %s\n", n, string(buff))
}

func (s *Server) newUID() (uint64, error) {
	for {
		n, err := rand.Int(rand.Reader, big.NewInt(math.MaxInt64))
		if err != nil {
			if err != nil {
				return 0, fmt.Errorf("failed to generate random number; %v", err)
			}
		}

		if _, ok := s.clientUids[n.String()]; !ok {
			s.clientUids[n.String()] = true
			return n.Uint64(), nil
		}

	}
}
