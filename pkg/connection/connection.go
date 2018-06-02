package connection

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"errors"
	"fmt"
	"io"
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
	sk   []byte
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

	k := d.CalcSharedSecret(new(big.Int).SetBytes(in))
	sha := sha256.New()
	if _, err := sha.Write(k.Bytes()); err != nil {
		return nil, fmt.Errorf("failed to sha256 shared secret: %v", err)
	}

	sk := sha.Sum(nil)

	return &Connection{
		conn: conn,
		dhke: d,
		sk:   sk,
	}, nil
}

func (c *Connection) Read() ([][]byte, error) {
	buff := make([]byte, 4096)

	n, err := c.conn.Read(buff)
	if err != nil {
		return nil, err
	}
	buff = buff[:n]

	buff, err = c.decrypt(buff)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt cipher: %v", err)
	}

	return c.decodeMessage(buff), nil
}

func (c *Connection) Write(b []byte) error {
	b, err := c.encypt(b)
	if err != nil {
		return fmt.Errorf("failed to encrypt message: %v", err)
	}

	if _, err := c.conn.Write(b); err != nil {
		return err
	}

	return nil
}

func (c *Connection) decodeMessage(d []byte) [][]byte {
	return bytes.Split(d, MessageBreak)
}

func (c *Connection) encypt(text []byte) ([]byte, error) {
	block, err := aes.NewCipher(c.sk)
	if err != nil {
		return nil, err
	}

	msg := c.pad(text)
	ciphertext := make([]byte, aes.BlockSize+len(msg))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil, err
	}

	cfb := cipher.NewCFBEncrypter(block, iv)
	cfb.XORKeyStream(ciphertext[aes.BlockSize:], []byte(msg))
	return ciphertext, nil
}

func (c *Connection) decrypt(text []byte) ([]byte, error) {
	block, err := aes.NewCipher(c.sk)
	if err != nil {
		return nil, err
	}

	if (len(text) % aes.BlockSize) != 0 {
		return nil, errors.New("blocksize must be multipe of decoded message length")
	}

	iv := text[:aes.BlockSize]
	msg := text[aes.BlockSize:]

	cfb := cipher.NewCFBDecrypter(block, iv)
	cfb.XORKeyStream(msg, msg)

	unpadMsg, err := c.unpad(msg)
	if err != nil {
		return nil, err
	}

	return unpadMsg, nil
}

func (c *Connection) pad(src []byte) []byte {
	padding := aes.BlockSize - len(src)%aes.BlockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(src, padtext...)
}

func (c *Connection) unpad(src []byte) ([]byte, error) {
	length := len(src)
	unpadding := int(src[length-1])

	if unpadding > length {
		return nil, errors.New("unpad error. This could happen when incorrect encryption key is used")
	}

	return src[:(length - unpadding)], nil
}
