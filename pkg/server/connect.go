package server

import (
	"crypto/rand"
	"crypto/x509"
	"fmt"
	"math/big"

	"github.com/joshvanl/go-whisper/pkg/connection"
)

const (
	MaxNumber = 99999999999
)

var (
	MessageBreak = []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
)

func (s *Server) Handle(conn *connection.Connection) {

	payload, err := conn.Read()
	if err != nil {
		return
	}

	if len(payload) == 0 {
		return
	}

	switch string(payload[0]) {
	case "first connection":
		if len(payload) != 3 {
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
		n, err := rand.Int(rand.Reader, big.NewInt(MaxNumber))
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

func (s *Server) newClient(conn *connection.Connection, recv [][]byte) error {
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
	if err = conn.Write(payload); err != nil {
		return fmt.Errorf("failed to send payload to client: %v", err)
	}

	return nil
}
