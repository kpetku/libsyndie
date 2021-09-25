package crypto

import (
	"crypto/sha256"

	"github.com/go-i2p/go-i2p/lib/common/base64"
	"github.com/go-i2p/go-i2p/lib/crypto"
)

// SigningKeypair is a DSA keypair used for replying to sign Syndie messages
type SigningKeypair struct {
	Pub  crypto.DSAPublicKey
	Priv crypto.DSAPrivateKey
}

// NewPrivateReplyKeypair creates a new DSA SigningKeypair
func NewSigningKeypair() *SigningKeypair {
	return new(SigningKeypair)
}

// Generate generates a new signing key pair
func (skp *SigningKeypair) Generate() error {
	privKey, err := skp.Priv.Generate()
	if err != nil {
		return err
	}
	skp.Priv = privKey
	pubKey, err := privKey.Public()
	if err != nil {
		return err
	}
	skp.Pub = pubKey
	return err
}

// String returns the base64 encoded public DSA key used for signing messages
func (i SigningKeypair) String() string {
	return base64.I2PEncoding.EncodeToString(i.Pub[:])
}

// Hash returns the sha256 base64 encoded short hash of the SigningKeypair public key
func (i SigningKeypair) Hash() string {
	foo := sha256.New()
	bar := base64.I2PEncoding.EncodeToString(i.Pub[:])
	foo.Write([]byte(bar))
	return base64.I2PEncoding.EncodeToString(foo.Sum(nil))
}
