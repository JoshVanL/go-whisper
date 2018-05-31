package key

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
)

const (
	keySize = 4056
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

func VerifyPayload(pk *rsa.PublicKey, payload string, sig []byte) error {
	opts := &rsa.PSSOptions{
		SaltLength: rsa.PSSSaltLengthEqualsHash,
		Hash:       crypto.SHA512,
	}

	hash := opts.Hash.New()
	_, err := hash.Write([]byte(payload))
	if err != nil {
		return fmt.Errorf("failed to hash payload: %v", err)
	}

	if err := rsa.VerifyPSS(pk, crypto.SHA512, hash.Sum(nil), []byte(sig), opts); err != nil {
		return fmt.Errorf("unable to verify payload: %v", err)
	}

	return nil
}

func SignMessage(pk *rsa.PrivateKey, message string) ([]byte, error) {
	opts := &rsa.PSSOptions{
		SaltLength: rsa.PSSSaltLengthEqualsHash,
		Hash:       crypto.SHA512,
	}

	hash := opts.Hash.New()
	_, err := hash.Write([]byte(message))
	if err != nil {
		return nil, fmt.Errorf("failed to hash message: %v", err)
	}
	hashed := hash.Sum(nil)

	signiture, err := pk.Sign(rand.Reader, hashed, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to sign message: %v", err)
	}

	return signiture, nil
}
