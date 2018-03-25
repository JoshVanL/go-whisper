package client

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"fmt"
	"net"

	"github.com/sirupsen/logrus"

	"github.com/joshvanl/go-whisper/pkg/key"
)

const (
	network = "tcp"
	keySize = 2048
)

type Client struct {
	log  *logrus.Entry
	addr string

	key  *rsa.PrivateKey
	conn net.Conn
	uid  uint64
}

func New(addr string, log *logrus.Entry) (*Client, error) {

	k, err := key.RetrieveLocalKey()
	if err != nil {
		return nil, fmt.Errorf("failed to read client key: %v", err)
	}

	uid, err := key.RetrieveLocalUID()
	if err != nil {
		return nil, fmt.Errorf("failed to read client config: %v", err)
	}

	return &Client{
		log: log,
		addr: addr,
		key: k,
		uid: uid,
	},
		nil
}

func (c *Client) Connect() error {
	conn, err := net.Dial(network, c.addr)
	if err != nil {
		return fmt.Errorf("failed to connect to server: %v", err)
	}

	c.conn = conn

	if _, err := c.conn.Write([]byte("hello")); err != nil {
		return fmt.Errorf("failed to write to connection: %v", err)
	}

	return nil
}

func (c *Client) encryptMessage(message, label []byte, k *rsa.PublicKey) ([]byte, error) {

	cipher, err := rsa.EncryptOAEP(sha256.New(), rand.Reader, k, message, label)
	if err != nil {
		return nil, fmt.Errorf("encryption of message failed: %v", err)
	}

	return cipher, nil
}

func (c *Client) decryptCipherText(cipher, label []byte) ([]byte, error) {

	plainText, err := rsa.DecryptOAEP(sha256.New(), rand.Reader, c.key, cipher, label)
	if err != nil {
		return nil, fmt.Errorf("decryption of ciphertext failed: %v", err)
	}

	return plainText, nil
}
