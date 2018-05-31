package client

import (
	//"crypto/x509"
	"fmt"

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

	fmt.Printf("%v\n", n)
	fmt.Printf("%v\n", d)

	return nil
}
