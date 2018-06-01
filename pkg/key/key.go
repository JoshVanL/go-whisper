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

type Key struct {
	dir string
	uid uint64
	sk  *rsa.PrivateKey
}

func New(dir string) (*Key, error) {
	key := &Key{
		dir: dir,
	}

	sk, err := key.retrieveLocalKey()
	if err != nil {
		return nil, err
	}
	key.sk = sk

	return key, nil
}

func (k *Key) createKeyPair() error {
	sk, err := k.generateRSAKey()
	if err != nil {
		return err
	}

	privBlock := &pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(sk)}
	pubBlock := &pem.Block{Type: "RSA PUBLIC KEY", Bytes: x509.MarshalPKCS1PublicKey(&sk.PublicKey)}

	if err := k.writeKeyPemFile(fmt.Sprintf("%s/%s", k.dir, privateKeyFile), privBlock); err != nil {
		return fmt.Errorf("failed to write private key to file: %v", err)
	}

	if err := k.writeKeyPemFile(fmt.Sprintf("%s/%s", k.dir, publicKeyFile), pubBlock); err != nil {
		return fmt.Errorf("failed to write public key to file: %v", err)
	}

	return nil
}

func (k *Key) generateRSAKey() (*rsa.PrivateKey, error) {
	sk, err := rsa.GenerateKey(rand.Reader, keySize)
	if err != nil {
		return nil, fmt.Errorf("failed to generate RSA key pair: %v", err)
	}

	return sk, nil
}

func (k *Key) VerifyPayload(pk *rsa.PublicKey, payload string, sig []byte) error {
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

func (k *Key) SignMessage(message string) ([]byte, error) {
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

	signiture, err := k.sk.Sign(rand.Reader, hashed, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to sign message: %v", err)
	}

	return signiture, nil
}

func (k *Key) Key() *rsa.PrivateKey {
	return k.sk
}

func (k *Key) Uid() uint64 {
	return k.uid
}
