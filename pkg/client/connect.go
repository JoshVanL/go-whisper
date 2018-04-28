package client

import (
	"crypto/x509"
	"fmt"
)

func (c *Client) Handshake() error {
	_, err := x509.MarshalPKIXPublicKey(c.key.PublicKey)
	if err != nil {
		return fmt.Errorf("failed to marshal public key: %v", err)
	}

	return nil
}
