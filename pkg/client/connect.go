package client

import (
	"crypto"
	"crypto/rsa"
	"crypto/x509"
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/joshvanl/go-whisper/pkg/key"
)

func (c *Client) Handshake() error {
	uids, err := key.NewUIDs(c.dir, c.config.UID)
	if err != nil {
		return err
	}
	c.uids = uids

	if c.uids.UID == 0 {
		return c.FirstConnection()
	}

	return nil
}

func (c *Client) FirstConnection() error {
	//d, err := x509.MarshalPKIXPublicKey(c.key.PublicKey)
	//if err != nil {
	//	return fmt.Errorf("failed to marshal public key: %v", err)
	//}

	_, err := c.conn.Write([]byte("first connection"))
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
	//fmt.Printf("%s\n", uidStr)
	//fmt.Printf("%s\n", pkStr)
	//fmt.Printf("%s\n", sigStr)

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

	if err := verifyPayload(pk, fmt.Sprintf("%s_%s", uidStr, pkStr), sigByte); err != nil {
		return err
	}

	return nil
}

func decodeMessage(d []byte) []string {
	return strings.Split(string(d), "_")
}

func verifyPayload(pk *rsa.PublicKey, payload string, sig []byte) error {
	opts := &rsa.PSSOptions{
		SaltLength: rsa.PSSSaltLengthEqualsHash,
		Hash:       crypto.SHA256,
	}

	hash := opts.Hash.New()
	_, err := hash.Write([]byte(payload))
	if err != nil {
		return fmt.Errorf("failed to hash payload: %v", err)
	}

	if err := rsa.VerifyPSS(pk, crypto.SHA256, hash.Sum(nil), []byte(sig), opts); err != nil {
		return fmt.Errorf("unable to verify payload from server: %v", err)
	}

	return nil
}
