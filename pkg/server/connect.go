package server

import (
	"fmt"
	"net"
	"crypto/rsa"
	"crypto/rand"
	"math"
	"math/big"
)


func (s *Server) Handle(con net.Conn) {
	buff := make([]byte, 256)

	_, err := con.Read(buff)
	if err != nil {
		return
	}

	uid, err := s.newUID()
	if err != nil {
		fmt.Errorf("failed to create new uid: %v", err)
	}

	fmt.Printf("uid: %v\n", uid)

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

		if _, ok := s.uids[n.Uint64()]; !ok {
			s.uids[n.Uint64()] = new(rsa.PublicKey)
			return n.Uint64(), nil
		}

	}
}
