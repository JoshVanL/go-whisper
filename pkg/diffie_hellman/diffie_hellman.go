package diffie_hellman

import (
	"crypto/rand"
	"errors"
	"fmt"
	"math/big"
)

const (
	P = "FFFFFFFFFFFFFFFFC90FDAA22168C234C4C6628B80DC1CD129024E088A67CC74020BBEA63B139B22514A08798E3404DDEF9519B3CD3A431B302B0A6DF25F14374FE1356D6D51C245E485B576625E7EC6F44C42E9A637ED6B0BFF5CB6F406B7EDEE386BFB5A899FA5AE9F24117C4B1FE649286651ECE45B3DC2007CB8A163BF0598DA48361C55D39A69163FA8FD24CF5F83655D23DCA3AD961C62F356208552BB9ED529077096966D670C354E4ABC9804F1746C08CA18217C32905E462E36CE3BE39E772C180E86039B2783A2EC07A28FB5C55DF06F4C52C9DE2BCBF6955817183995497CEA956AE515D2261898FA051015728E5A8AACAA68FFFFFFFFFFFFFFFF"
	G = 3
)

type DiffieHellman struct {
	s, p, g *big.Int
	k, in   *big.Int
}

func New() (*DiffieHellman, error) {

	g := big.NewInt(G)
	p, succ := new(big.Int).SetString(P, 16)
	if !succ {
		return nil, errors.New("faild to set p prime value")
	}

	max := new(big.Int)
	max.Exp(big.NewInt(2), big.NewInt(2048), nil).Sub(max, big.NewInt(1))

	s, err := rand.Int(rand.Reader, max)
	if err != nil {
		return nil, fmt.Errorf("failed to generate secret: %v", err)
	}

	return &DiffieHellman{
		g:  g,
		p:  p,
		s:  s,
		in: new(big.Int).Exp(g, s, p),
	}, nil
}

func (d *DiffieHellman) Intermediate() *big.Int {
	return d.in
}

func (d *DiffieHellman) CalcSharedSecret(b *big.Int) *big.Int {
	d.k = new(big.Int).Exp(b, d.s, d.p)
	return d.k
}

func (d *DiffieHellman) SharedSecret() *big.Int {
	return d.k
}
