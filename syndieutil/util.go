package syndieutil

import (
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"io"
	"strings"

	"github.com/go-i2p/go-i2p/lib/common/base64"
)

const newLine string = "\n"

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

func readInt(r io.Reader) int {
	buf := make([]byte, 4)
	io.ReadFull(r, buf)
	internalSize := binary.BigEndian.Uint32(buf)
	return int(internalSize)
}

func value(s string) (string, error) {
	if strings.Contains(s, "=") {
		return strings.Join(strings.SplitAfter(strings.TrimSpace(s), "=")[1:], ""), nil
	}
	return "", fmt.Errorf("invalid string: %s", s)
}
