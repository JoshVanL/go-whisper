package client

import (
	"crypto/rsa"
	"fmt"
	"net"
	"sort"

	"github.com/joshvanl/go-whisper/pkg/config"
	"github.com/joshvanl/go-whisper/pkg/connection"
	"github.com/joshvanl/go-whisper/pkg/gui"
	"github.com/joshvanl/go-whisper/pkg/key"
)

const (
	network = "tcp"
	keySize = 2048
)

type Client struct {
	addr string
	dir  string

	key      *key.Key
	conn     *connection.Connection
	serverpk *rsa.PublicKey

	config *config.Config
	g      *gui.GUI
}

func New(addr, dir string) (*Client, error) {

	g, err := gui.New()
	if err != nil {
		return nil, fmt.Errorf("failed to initiate gui: %v", err)
	}

	g.Infof("Retrieving local key pair...")
	k, err := key.New(dir)
	if err != nil {
		return nil, fmt.Errorf("failed to read client key: %v", err)
	}

	client := &Client{
		dir: dir,
		key: k,
		g:   g,
	}

	g.SetClient(client)

	g.Infof("Retrieving local client config...")
	config, err := config.ReadConfig(dir)
	if err != nil {
		return nil, fmt.Errorf("failed to read config: %v", err)
	}
	client.config = config
	client.addr = client.config.Address

	if addr != "" {
		client.addr = addr
	}

	g.Infof("Connecting to server...")

	return client, nil
}

func (c *Client) Close() {
	if c != nil && c.g != nil {
		c.g.Close()
	}
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

	c.g.Infof("Connection successful.")

	c.g.SetUid(c.config.UID)

	c.g.DrawMenu()

	return nil
}

func (c *Client) Uids() []string {
	mapUids, err := c.key.UIDsFromFile()
	if err != nil {
		// do something here!
		return nil
	}

	var uids []string
	for key, b := range mapUids {
		if b && key != "0" {
			uids = append(uids, key)
		}
	}

	sort.Strings(uids)

	return uids
}
