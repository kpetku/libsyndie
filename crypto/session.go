package crypto

import (
	"crypto/rand"

	"github.com/go-i2p/go-i2p/lib/common/base64"
)

// SessionKey is a base64 encoded AES-256 key
type SessionKey = string

// NewSessionKey creates a new base64 encoded AES-256 key
func NewSessionKey() SessionKey {
	key := make([]byte, 32)
	rand.Read(key)
	return base64.I2PEncoding.EncodeToString(key)
}
