package server

import (
	"bytes"
	"crypto/rand"
	"crypto/x509"
	"fmt"
	"math"
	"math/big"
	"net"
)

var (
	MessageBreak = []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
)

func (s *Server) Handle(conn net.Conn) {
	buff := make([]byte, 4096)

	n, err := conn.Read(buff)
	if err != nil {
		return
	}
	buff = buff[:n]

	fmt.Printf("buff: %s\n", buff[:n])

	payload := decodeMessage(buff)

	if len(payload) == 0 {
		return
	}

	if string(payload[0]) == "first connection" {
		if len(payload) != 2 {
			return
		}

		err = s.newClient(conn, payload)
		if err != nil {
			s.log.Errorf("error handling new user: %v", err)
		}

		return
	}
}

func (s *Server) newUID() (string, error) {
	for {
		n, err := rand.Int(rand.Reader, big.NewInt(math.MaxInt64))
		if err != nil {
			if err != nil {
				return "", fmt.Errorf("failed to generate random number; %v", err)
			}
		}

		if _, ok := s.clientUids[n.String()]; !ok {
			s.clientUids[n.String()] = true
			return n.String(), nil
		}

	}
}

func (s *Server) newClient(conn net.Conn, recv [][]byte) error {
	uid, err := s.newUID()
	if err != nil {
		fmt.Errorf("failed to create new uid: %v", err)
	}

	pk, err := x509.ParsePKCS1PublicKey(recv[1])
	if err != nil {
		return fmt.Errorf("failed to parse client public key: %v", err)
	}

	if err := s.key.NewUidFile(uid, pk); err != nil {
		return fmt.Errorf("failed to store client public key: %v", err)
	}

	message := append(append([]byte(uid), MessageBreak...), s.key.PublicKey()...)
	signiture, err := s.key.SignMessage(message)
	if err != nil {
		return fmt.Errorf("failed to sign message for client: %v", err)
	}

	payload := append(append(message, MessageBreak...), signiture...)
	_, err = conn.Write([]byte(payload))
	if err != nil {
		return fmt.Errorf("failed to send payload to client: %v", err)
	}

	return nil
}

func decodeMessage(d []byte) [][]byte {
	return bytes.Split(d, MessageBreak)
}
