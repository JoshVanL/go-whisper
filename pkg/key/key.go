package key

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
)

const (
	keySize = 2048
)

func CreateKeyPair(dir string) error {
	k, err := generateRSAKey()
	if err != nil {
		return err
	}

	privBlock := &pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(k)}
	pubBlock := &pem.Block{Type: "RSA PUBLIC KEY", Bytes: x509.MarshalPKCS1PublicKey(&k.PublicKey)}

	if err := writeKeyPemFile(fmt.Sprintf("%s/%s", dir, privateKeyFile), privBlock); err != nil {
		return fmt.Errorf("failed to write private key to file: %v", err)
	}

	if err := writeKeyPemFile(fmt.Sprintf("%s/%s", dir, publicKeyFile), pubBlock); err != nil {
		return fmt.Errorf("failed to write public key to file: %v", err)
	}

	return nil
}

func generateRSAKey() (*rsa.PrivateKey, error) {
	k, err := rsa.GenerateKey(rand.Reader, keySize)
	if err != nil {
		return nil, fmt.Errorf("failed to generate RSA key pair: %v", err)
	}

	return k, nil
}
