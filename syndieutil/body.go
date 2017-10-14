package syndieutil

import (
	"archive/zip"
	"bufio"
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/hkparker/go-i2p/lib/common/base64"
)

func (header *Header) Unmarshal(r io.Reader) (*Message, error) {
	var rest bytes.Buffer

	br := bufio.NewReader(r)
	var realSize int

	line, lerr := br.ReadString('\n')
	if lerr != nil {
		return nil, errors.New("invalid message: " + lerr.Error())
	}
	// find the magic "Syndie.Message.1." string
	if !strings.HasPrefix(line, "Syndie.Message.1.") {
		return nil, errors.New("invalid message")
	}
	for {
		line, err := br.ReadString('\n')
		if err != nil {
			break
		}
		if line == "\n" {
			size, err := br.ReadString('\n')
			if err != nil {
				break
			}
			foo, err := value(size)
			bar, err := strconv.Atoi(foo)
			if err != nil {
				break
			}
			realSize = bar
			break
		}
		// call ReadLine on each line
		err = header.ReadLine(strings.TrimSpace(line))
		if err != nil {
			return nil, errors.New("error from header ReadLine: " + line + " error: " + err.Error())
		}
	}

	// read the enclosed message body into enclosed
	var enclosed = make([]byte, realSize)
	_, err := io.ReadFull(br, enclosed)
	if err != nil {
		return nil, errors.New("error in ioReadFull: " + err.Error())
	}
	rest.Write(enclosed)

	if len(enclosed) < 32 {
		return nil, errors.New("invalid message: too small for IV")
	}

	decrypted := make([]byte, len(enclosed)+32)

	// taken from the raw, encrypted enclosed message
	iv := enclosed[0:16]

	// after this point, the stuff below needs to be decrypted!
	k, err := base64.I2PEncoding.DecodeString(header.BodyKey)

	if err != nil {
		return nil, errors.New("error decoding: " + err.Error())
	}

	block, err := aes.NewCipher([]byte(k))
	if err != nil {
		return nil, errors.New("error initializing NewCipher: %s" + err.Error())
	}

	decrypter := cipher.NewCBCDecrypter(block, iv)
	decrypter.CryptBlocks(decrypted, enclosed[16:realSize])

	start := bytes.IndexByte(decrypted, 0x0)

	internalSize := binary.BigEndian.Uint32(decrypted[start+1 : start+5])
	totalSize := binary.BigEndian.Uint32(decrypted[start+5 : start+9])

	if realSize != int(totalSize)+16 {
		return nil, fmt.Errorf("size mismatch, expected: %d, found: %d", realSize, int(totalSize)+16)
	}

	zr, err := zip.NewReader(bytes.NewReader(decrypted[start+9:start+9+int(internalSize)]), int64(start+9-start+9+int(internalSize)))

	if err != nil {
		return nil, err
	}

	// hand off the decrypted zip to ParseMessage
	m, err := header.ParseMessage(zr)
	if err != nil {
		return nil, err
	}
	// reached the end of the body, next comes the signature area
	authorizationSig, err := br.ReadString('\n')
	rest.Write([]byte(authorizationSig))
	if err != nil {
		return nil, err
	}
	//	foo, err := value(authorizationSig)
	//	log.Printf("authorizationSig: %s", foo)

	authenticationSig, err := br.ReadString('\n')
	rest.Write([]byte(authenticationSig))

	if err != nil {
		return nil, err
	}
	//	bar, err := value(authenticationSig)
	//	log.Printf("authenticationSig: %s", bar)

	// TODO: lots

	// check the hmac
	var hmacPreKey bytes.Buffer
	hmacPreKey.Write(k)
	hmacPreKey.Write(iv)

	sha := sha256.New()
	sha.Write(hmacPreKey.Bytes())

	h := hmac.New(sha256.New, sha.Sum(nil))
	h.Write(rest.Bytes()[16 : realSize-32])

	if !hmac.Equal(h.Sum(nil), rest.Bytes()[realSize-32:realSize]) {
		return nil, fmt.Errorf("unable to verify HMAC")
	}
	return &m, nil
}

func value(s string) (string, error) {
	if strings.Contains(s, "=") {
		return strings.Join(strings.SplitAfter(strings.TrimSpace(s), "=")[1:], ""), nil
	}
	return "", fmt.Errorf("invalid string: %s", s)
}
