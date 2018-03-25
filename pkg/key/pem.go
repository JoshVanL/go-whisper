package key

import (
	"fmt"
	"os/user"
	"os"
	"encoding/pem"
	"crypto/x509"
	"crypto/rsa"
	"io/ioutil"
)

const (
	wipserDir      = ".go-wipser"
	uidDir         = "uid_public_keys"
	privateKeyFile = "private_key.pem"
	publicKeyFile  = "public_key.pem"
)

func RetrieveLocalKey() (*rsa.PrivateKey, error) {
	if err := ensureKeyDirectory(); err != nil {
		return nil, fmt.Errorf("failed to ensure key directory: %v", err)
	}

	if err := ensureKeyFiles(); err != nil {
		return nil, fmt.Errorf("failed to ensure key files: %v", err)
	}

	dir, err := keyDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get key directory: %v", err)
	}

	k, err := readPrivateKey(fmt.Sprintf("%s/%s", dir, privateKeyFile))
	if err != nil {
		return nil, fmt.Errorf("failed to read key: %v", err)
	}

	return k, nil
}

func readPrivateKey(path string) (*rsa.PrivateKey, error) {
	f, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open private key file for reading: %v", err)
	}

	block, rest := pem.Decode(f)
	if block == nil {
		return nil, fmt.Errorf("pem block was nil at private key file: %s", path)
	}
	if len(rest) != 0 {
		return nil, fmt.Errorf("expected rest of pem block to be nil, got=%v", rest)
	}

	k, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse private key file: %v", err)
	}

	return k, nil
}

func readPublicKey(path string) (*rsa.PublicKey, error) {
	f, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open public key file for reading: %v", err)
	}

	block, rest := pem.Decode(f)
	if block == nil {
		return nil, fmt.Errorf("pem block was nil at public key file: %s", path)
	}
	if len(rest) != 0 {
		return nil, fmt.Errorf("expected rest of pem block to be nil, got=%v", rest)
	}

	k, err := x509.ParsePKCS1PublicKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse public key file: %v", err)
	}

	return k, nil
}


func ensureKeyDirectory() error {
	dir, err := keyDir()
	if err != nil {
		return fmt.Errorf("failed to get key directory: %v", err)
	}

	if _, err := os.Stat(dir); err != nil {
		if os.IsNotExist(err) {
			if err := os.MkdirAll(dir, 0700); err != nil {
				return fmt.Errorf("failed to create go-wisper key directory: %v", err)
			}

		} else {
			return fmt.Errorf("failed to check go-wisper key directory: %v", err)
		}
	}

	return nil
}

func uidPubkicKeyDirectory() (string, error) {
	dir, err := uidPublicKeyDir()
	if err != nil {
		return "", fmt.Errorf("failed to get uids public key directory: %v", err)
	}

	if f, err := os.Stat(dir); err != nil {
		if os.IsNotExist(err) {
			if err := os.MkdirAll(dir, 0700); err != nil {
				return "", fmt.Errorf("failed to create uids public keys directory: %v", err)
			}

		} else {
			return "", fmt.Errorf("failed to check uids public keys directory: %v", err)
		}
	} else {
		if !f.IsDir() {
			return "", fmt.Errorf("is not a directory: %s", dir)
		}
	}

	return dir, nil
}

func ensureKeyFiles() error {
	dir, err := keyDir()
	if err != nil {
		return fmt.Errorf("failed to get key directory: %v", err)
	}

	if _, err := os.Stat(fmt.Sprintf("%s/%s", dir, privateKeyFile)); err != nil {
		if os.IsNotExist(err) {
			if err := CreateKeyPair(dir); err != nil {
				return fmt.Errorf("failed to create new key pair: %v", err)
			}
		} else {
			return fmt.Errorf("error checking key files: %v", err)
		}
	}

	return nil
}


func keyDir() (string, error) {
	usr, err := user.Current()
	if err != nil {
		return "", fmt.Errorf("unable to resolve current user: %v", err)
	}

	_, err = os.Stat(usr.HomeDir)
	if err != nil {
		return "", fmt.Errorf("error checking home directory : %v", err)
	}

	return fmt.Sprintf("%s/%s", usr.HomeDir, wipserDir), nil
}

func uidPublicKeyDir() (string, error) {
	usr, err := user.Current()
	if err != nil {
		return "", fmt.Errorf("unable to resolve current user: %v", err)
	}

	_, err = os.Stat(usr.HomeDir)
	if err != nil {
		return "", fmt.Errorf("error checking home directory : %v", err)
	}

	return fmt.Sprintf("%s/%s/%s", usr.HomeDir, wipserDir, uidPublicKeyDir), nil
}


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

func writeKeyPemFile(file string, key *pem.Block) error {
	f, err := os.OpenFile(file, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return fmt.Errorf("failed to open pem file: %v", err)
	}
	defer f.Close()

	if err := pem.Encode(f, key); err != nil {
		return fmt.Errorf("failed to write pem block to file: %v", err)
	}

	return nil
}
