package enclosure

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/binary"

	"github.com/kpetku/go-syndie/lib/common"

	log "github.com/Sirupsen/logrus"
)

// SyndieTrailer contains a Syndie trailer that contains version and pairs fields
type SyndieTrailer struct {
	Size              int
	Raw               []byte
	decrypted         []byte
	AuthorizationSig  []byte
	AuthenticationSig []byte
	Body              []byte
	EOF               int
}

// SyndieTrailerPayload holds the following: rand(nonzero) padding + 0 + internalSize + totalSize + data + rand
type SyndieTrailerPayload struct {
	InternalSize int
	TotalSize    int

	Decrypted    []byte
	DecryptedRaw []byte
	IV           []byte
	HMAC         []byte
	BodySection  []byte
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
	inner.IV = trailer.Raw[0:16]
	inner.HMAC = trailer.Raw[trailer.Size-32 : trailer.Size]

	decrypter := cipher.NewCBCDecrypter(block, trailer.Raw[:16])
	decrypted := make([]byte, trailer.Size+32)
	decrypter.CryptBlocks(decrypted, trailer.Raw[16:trailer.Size])
	for i := range decrypted {
		if !found && decrypted[i] == 0x0 {
			is := int(binary.BigEndian.Uint32(decrypted[i+1 : i+5]))
			ts := int(binary.BigEndian.Uint32(decrypted[i+5 : i+9]))
			if trailer.Size != ts+16 {
				log.WithFields(log.Fields{
					"at":           "(trailer) DecryptAES",
					"trailer_size": trailer.Size,
					"is":           is,
					"ts":           ts + 16,
					"reason":       "payload size did not match envelope size",
				}).Fatalf("%s", err)
			}
			inner.InternalSize = is
			inner.TotalSize = ts
			inner.Decrypted = decrypted[i+9 : i+9+is]
			inner.DecryptedRaw = decrypted

			var hmacPreKey bytes.Buffer
			hmacPreKey.Write(k)
			hmacPreKey.Write(trailer.Raw[0:16])

			sha := sha256.New()
			sha.Write(hmacPreKey.Bytes())

			if !verifyHmac256(trailer.Raw[16:trailer.Size-32], trailer.Raw[trailer.Size-32:trailer.Size], sha.Sum(nil)) {
				log.WithFields(log.Fields{
					"at":     "(trailer) DecryptAES, verifyHmac256",
					"reason": "invalid HMAC",
				}).Fatalf("%s", err)
			}
			found = true
		}
	}

	return inner
}

func verifyHmac256(stringToVerify []byte, signature []byte, sharedSecret []byte) bool {
	h := hmac.New(sha256.New, sharedSecret)
	h.Write(stringToVerify)
	calculated := h.Sum(nil)
	return hmac.Equal(calculated, signature)
}

func colonDelimiter(data []byte, atEOF bool) (advance int, token []byte, err error) {
	if atEOF && len(data) == 0 {
		return 0, nil, nil
	}
	if i := bytes.Index(data, []byte{58}); i >= 0 {
		return i + 1, data[0:i], nil
	}
	if atEOF {
		return len(data), data, nil
	}
	return 0, nil, nil
}
