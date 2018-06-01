package client

import (
	"crypto/x509"
	"encoding/hex"
	"fmt"
	"strings"
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
	_, err := c.conn.Write(append([]byte("first connection_"), []byte(c.key.PublicKey())...))
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

	pkByte, err := hex.DecodeString(string(pkStr))
	if err != nil {
		return fmt.Errorf("failed to decode public key hex string: %v", err)
	}

	sigByte, err := hex.DecodeString(string(sigStr))
	if err != nil {
		fmt.Errorf("failed to decode signiture hex string: %v", err)
	}

	pk, err := x509.ParsePKCS1PublicKey(pkByte)
	if err != nil {
		return fmt.Errorf("failed to parse server public key: %v", err)
	}

	if err := c.key.VerifyPayload(pk, fmt.Sprintf("%s_%s", uidStr, pkStr), sigByte); err != nil {
		return err
	}

	if err := c.key.NewUidFile("0", pk); err != nil {
		return err
	}

	return nil
}

func decodeMessage(d []byte) []string {
	return strings.Split(string(d), "_")
}
