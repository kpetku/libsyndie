package enclosure

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/binary"

	log "github.com/Sirupsen/logrus"
	"github.com/kpetku/go-syndie/lib/common"
)

// SyndieTrailer contains a Syndie trailer that contains version and pairs fields
type SyndieTrailer struct {
	size              int
	raw               []byte
	authorizationSig  string
	authenticationSig string
}

// SyndieTrailerPayload holds the following: rand(nonzero) padding + 0 + internalSize + totalSize + data + rand
type SyndieTrailerPayload struct {
	internalSize int
	totalSize    int
	decrypted    []byte
}

// DecryptAES decrypts a SyndieTrailer into a messagePayload
func (trailer *SyndieTrailer) DecryptAES(key string) SyndieTrailerPayload {
	var found bool
	inner := SyndieTrailerPayload{}
	k, err := common.I2PEncoding.DecodeString(key)
	if err != nil {
		log.WithFields(log.Fields{
			"at":     "(trailer) DecryptAES, DecodeString",
			"key":    key,
			"reason": "unable to convert key into b64",
		}).Fatalf("%s", err)
	}
	block, err := aes.NewCipher([]byte(k))
	if err != nil {
		log.WithFields(log.Fields{
			"at":     "(trailer) DecryptAES, NewCipher",
			"key":    key,
			"block":  block,
			"reason": "invalid block AES cipher",
		}).Fatalf("%s", err)
	}
	decrypter := cipher.NewCBCDecrypter(block, trailer.raw[:16])
	decrypted := make([]byte, len(trailer.raw[:trailer.size]))
	decrypter.CryptBlocks(decrypted, trailer.raw[:trailer.size])
	for i := range decrypted {
		if !found && decrypted[i] == 0x0 {
			is := int(binary.BigEndian.Uint32(decrypted[i+1 : +i+5]))
			ts := int(binary.BigEndian.Uint32(decrypted[i+5 : +i+9]))
			if trailer.size != ts+16 {
				log.WithFields(log.Fields{
					"at":           "(trailer) DecryptAES",
					"trailer_size": trailer.size,
					"is":           is,
					"ts":           ts + 16,
					"reason":       "payload size did not match envelope size",
				}).Fatalf("%s", err)
			}

			inner.internalSize = is
			inner.totalSize = ts
			inner.decrypted = decrypted[i+10 : is]
			found = true
		}
	}
	return inner
}
