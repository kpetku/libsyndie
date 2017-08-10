package syndieutil

import (
	"crypto/sha256"

	"github.com/hkparker/go-i2p/lib/common/base64"
)

func ShortIdent(i string) string {
	if len(i) > 6 {
		return "[" + i[0:6] + "]"
	}
	return "[" + i + "]"
}

func ChanHash(s string) (string, error) {
	foo := sha256.New()
	bar, err := base64.I2PEncoding.DecodeString(s)
	if err != nil {
		return "", err
	}
	foo.Write([]byte(bar))
	if err != nil {
		return "", err
	}
	return base64.I2PEncoding.EncodeToString(foo.Sum(nil)), nil
}
