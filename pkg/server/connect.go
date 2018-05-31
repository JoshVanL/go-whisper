package server

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/hex"
	"fmt"
	"math"
	"math/big"
	"net"
)

func (s *Server) Handle(conn net.Conn) {
	buff := make([]byte, 4096)

	n, err := conn.Read(buff)
	if err != nil {
		return
	}
	buff = buff[:n]

	fmt.Printf("buff: %s\n", buff[:n])

	if string(buff) == "first connection" {
		err = s.newClient(conn)
		if err != nil {
			s.log.Errorf("error handling new user: %v", err)
		}

		return
	}

	return

	//fmt.Printf("(%d) read from connection: %s\n", n, string(buff))
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

func (s *Server) newClient(conn net.Conn) error {
	uid, err := s.newUID()
	if err != nil {
		fmt.Errorf("failed to create new uid: %v", err)
	}

	d := x509.MarshalPKCS1PublicKey(&s.key.PublicKey)
	messge := fmt.Sprintf("%s____%s", uid, hex.EncodeToString(d))

	opts := &rsa.PSSOptions{
		SaltLength: rsa.PSSSaltLengthEqualsHash,
		Hash:       crypto.SHA256,
	}

	hash := opts.Hash.New()
	_, err = hash.Write([]byte(messge))
	if err != nil {
		return fmt.Errorf("failed to hash message: %v", err)
	}
	hashed := hash.Sum(nil)

	signiture, err := s.key.Sign(rand.Reader, hashed, opts)
	if err != nil {
		return fmt.Errorf("failed to sign message for client: %v", err)
	}

	payload := fmt.Sprintf("%s____%s", messge, hex.EncodeToString(signiture))
	_, err = conn.Write([]byte(payload))
	if err != nil {
		return fmt.Errorf("failed to send payload to client: %v", err)
	}

	return nil
}
