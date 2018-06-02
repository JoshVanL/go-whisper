package client

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"fmt"
	"net"

	"github.com/sirupsen/logrus"

	"github.com/joshvanl/go-whisper/pkg/config"
	"github.com/joshvanl/go-whisper/pkg/connection"
	"github.com/joshvanl/go-whisper/pkg/key"
)

const (
	network = "tcp"
	keySize = 2048
)

type Client struct {
	log *logrus.Entry

	addr string
	dir  string

	key  *key.Key
	conn *connection.Connection

	config *config.Config
}

func New(addr, dir string, log *logrus.Entry) (*Client, error) {

	log.Infof("Retrieving local key pair...")
	k, err := key.New(dir)
	if err != nil {
		return nil, fmt.Errorf("failed to read client key: %v", err)
	}

	client := &Client{
		log: log,
		dir: dir,
		key: k,
	}

	log.Infof("Retrieving local client config...")
	config, err := config.ReadConfig(dir)
	if err != nil {
		return nil, fmt.Errorf("failed to read config: %v", err)
	}
	client.config = config
	client.addr = client.config.Address

	if addr != "" {
		client.addr = addr
	}

	log.Infof("Connecting to server...")
	if err := client.Connect(); err != nil {
		return nil, err
	}

	return client, nil
}

func (c *Client) Connect() error {
	conn, err := net.Dial(network, c.addr)
	if err != nil {
		return fmt.Errorf("failed to connect to server: %v", err)
	}
	c.conn, err = connection.New(conn)
	if err != nil {
		return err
	}

	if err := c.Handshake(); err != nil {
		return fmt.Errorf("failed to handshake with the server: %v", err)
	}

	c.log.Infof("Connection successful.")

	//if _, err := c.conn.Write([]byte("hello")); err != nil {
	//	return fmt.Errorf("failed to write to connection: %v", err)
	//}

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

	plainText, err := rsa.DecryptOAEP(sha256.New(), rand.Reader, c.key.Key(), cipher, label)
	if err != nil {
		return nil, fmt.Errorf("decryption of ciphertext failed: %v", err)
	}

	return plainText, nil
}
