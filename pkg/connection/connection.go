package connection

import (
	"bytes"
	"fmt"
	"math/big"
	"net"

	dhke "github.com/joshvanl/go-whisper/pkg/diffie_hellman"
)

var (
	MessageBreak = []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
)

type Connection struct {
	conn net.Conn
	dhke *dhke.DiffieHellman
}

func New(conn net.Conn) (*Connection, error) {
	d, err := dhke.New()
	if err != nil {
		return nil, fmt.Errorf("failed to init Diffie Hellman: %v", err)
	}

	if _, err := conn.Write(d.Intermediate().Bytes()); err != nil {
		return nil, fmt.Errorf("failed to send Diffie Hellman intermediate: %v", err)
	}

	in := make([]byte, 2048)
	n, err := conn.Read(in)
	if err != nil {
		return nil, fmt.Errorf("failed to read Diffie Hellman intermediate: %v", err)
	}
	in = in[:n]

	d.CalcSharedSecret(new(big.Int).SetBytes(in))

	return &Connection{
		conn: conn,
		dhke: d,
	}, nil
}

func (c *Connection) Read() ([][]byte, error) {
	buff := make([]byte, 4096)

	n, err := c.conn.Read(buff)
	if err != nil {
		return nil, err
	}
	buff = buff[:n]

	return c.decodeMessage(buff), nil
}

func (c *Connection) Write(b []byte) error {
	if _, err := c.conn.Write(b); err != nil {
		return err
	}

	return nil
}

func (c *Connection) decodeMessage(d []byte) [][]byte {
	return bytes.Split(d, MessageBreak)
}
