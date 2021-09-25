package crypto

import (
	"crypto/rand"

	"github.com/go-i2p/go-i2p/lib/common/base64"
	"github.com/go-i2p/go-i2p/lib/crypto"
	"golang.org/x/crypto/openpgp/elgamal"
)

// PrivateReplyKeypair is an ElGamal keypair used for replying to private Syndie messages
type PrivateReplyKeypair struct {
	PubKey  elgamal.PublicKey
	PrivKey elgamal.PrivateKey
}

// NewPrivateReplyKeypair creates a new ElGamal PrivateReplyKeypair
func NewPrivateReplyKeypair() (*PrivateReplyKeypair, error) {
	priv := elgamal.PrivateKey{}
	prk := PrivateReplyKeypair{}
	err := crypto.ElgamalGenerate(&priv, rand.Reader)
	if err != nil {
		return &PrivateReplyKeypair{}, err
	}
	prk.PrivKey = priv
	prk.PubKey = priv.PublicKey
	return &prk, err
}

// String returns the base64 encoded public elgamal key used for private message replies
func (r PrivateReplyKeypair) String() string {
	// TODO: r.PrivKey.Y.Bytes() ...?
	return base64.I2PEncoding.EncodeToString(r.PrivKey.Y.Bytes())
}
