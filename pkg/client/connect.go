package client

import (
	"bytes"
	"crypto/x509"
	"fmt"
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

	return nil
}

func (c *Client) FirstConnection() error {
	_, err := c.conn.Write(append(append([]byte("first connection"), MessageBreak...), []byte(c.key.PublicKey())...))
	if err != nil {
		return fmt.Errorf("failed to write first connection: %v", err)
	}

	d := make([]byte, 4096)
	n, err := c.conn.Read(d)
	if err != nil {
		return fmt.Errorf("failed to read from connection: %v", err)
	}
	d = d[:n]

	rec := decodeMessage(d)
	if len(rec) != 3 {
		return fmt.Errorf("unexpected number of response, exp=3 got=%d", len(rec))
	}
	uidStr, pkStr, sigStr := rec[0], rec[1], rec[2]

	pk, err := x509.ParsePKCS1PublicKey(pkStr)
	if err != nil {
		return fmt.Errorf("failed to parse server public key: %v", err)
	}
	if err := c.key.VerifyPayload(pk, append(append(uidStr, MessageBreak...), pkStr...), sigStr); err != nil {
		return err
	}

	if err := c.key.NewUidFile("0", pk); err != nil {
		return err
	}

	return nil
}

func decodeMessage(d []byte) [][]byte {
	return bytes.Split(d, MessageBreak)
}
