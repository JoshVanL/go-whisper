package client

import (
	"crypto/x509"
	"errors"
	"fmt"
	"strconv"
)

var (
	MessageBreak = []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
)

func (c *Client) Handshake() error {
	if err := c.key.NewUIDs(c.config.UID); err != nil {
		return err
	}

	if c.key.Uid() == 0 {
		return c.FirstConnection()
	}

	pk, err := c.key.ReadUidFile("0")
	if err != nil {
		return fmt.Errorf("failed to read server public key from file: %v", err)
	}

	c.serverpk = pk

	return nil
}

func (c *Client) FirstConnection() error {
	send := appendParams([]byte("first connection"), c.key.PublicKey())

	signiture, err := c.key.SignMessage(send)
	if err != nil {
		return fmt.Errorf("failed to sign initial message: %v", err)
	}

	send = appendParams(send, signiture)
	if err = c.conn.Write(send); err != nil {
		return fmt.Errorf("failed to write first connection: %v", err)
	}

	rec, payload, err := c.conn.Read()
	if err != nil {
		return fmt.Errorf("failed to read from connection: %v", err)
	}

	if len(rec) != 3 {
		return fmt.Errorf("unexpected number of response, exp=3 got=%d", len(rec))
	}
	uidB, pkB, sigB := rec[0], rec[1], rec[2]

	pk, err := x509.ParsePKCS1PublicKey(pkB)
	if err != nil {
		return fmt.Errorf("failed to parse server public key: %v", err)
	}
	if err := c.key.VerifyPayload(pk, payload, sigB); err != nil {
		return err
	}

	c.serverpk = pk

	c.config.UID, err = strconv.ParseUint(string(uidB), 10, 64)
	if err != nil {
		return fmt.Errorf("failed to convert uid string to uint64: %v", err)
	}

	if err := c.config.Write(); err != nil {
		return err
	}

	if err := c.key.NewUidFile("0", pk); err != nil {
		return err
	}

	return nil
}

func (c *Client) QueryUID(uid string) (string, error) {

	message := appendParams([]byte("uid query"), []byte(fmt.Sprintf("%v", c.config.UID)))
	message = appendParams(message, []byte(uid))
	signiture, err := c.key.SignMessage(message)
	if err != nil {
		return "", fmt.Errorf("failed to sign query message: %v", err)
	}
	message = appendParams(message, signiture)

	if err := c.conn.Write(message); err != nil {
		return "", fmt.Errorf("failed to send uid query: %v", err)
	}

	res, payload, err := c.conn.Read()
	if err != nil {
		return "", err
	}

	if err := c.key.VerifyPayload(c.serverpk, payload, res[len(res)-1]); err != nil {
		return "", err
	}

	if len(res) < 3 {
		return "", errors.New(string(res[0]))
	}

	pk, err := x509.ParsePKCS1PublicKey(res[2])
	if err != nil {
		return "", fmt.Errorf("failed to parse uid public key: %v", err)
	}

	if err := c.key.NewUidFile(uid, pk); err != nil {
		return "", fmt.Errorf("failed to save new uid public key: %v", err)
	}

	return string(res[0]), nil
}

func appendParams(a, b []byte) []byte {
	return append(append(a, MessageBreak...), b...)
}
