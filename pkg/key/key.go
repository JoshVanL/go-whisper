package key

import (
	"crypto/rand"
	"crypto/rsa"
	"fmt"
)

const (
	keySize = 2048
)

func generateRSAKey() (*rsa.PrivateKey, error) {
	k, err := rsa.GenerateKey(rand.Reader, keySize)
	if err != nil {
		return nil, fmt.Errorf("failed to generate RSA key pair: %v", err)
	}

	return k, nil
}
