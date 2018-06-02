package connection

import (
	"bytes"
	"net"
)

var (
	MessageBreak = []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
)

type Connection struct {
	conn net.Conn
}

func New(conn net.Conn) *Connection {
	return &Connection{
		conn: conn,
	}
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
