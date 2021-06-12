package syndieutil

import (
	"archive/zip"
	"bufio"
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/sha256"
	"errors"
	"io"
	"strconv"
	"strings"

	"github.com/go-i2p/go-i2p/lib/common/base64"
)

const syndieMessage = "Syndie.Message.1."
const invalidMessage = "invalid message"
const limit = 1024

func (h *Header) Unmarshal(r io.Reader) (*Message, error) {
	h.reader = bufio.NewReader(r)
	for state := 0; state < int(invalid); state++ {
		err := h.next()
		if err != nil || h.err != nil {
			return nil, err
		}
		h.state++
	}
	return h.msg, nil
}

func (h *Header) next() (err error) {
	switch h.state {
	case readMagicVersionLine:
		h.err = h.readMagicVersionLine()
	case readHeaderKeyPairs:
		h.err = h.readHeaderKeyPairs()
	case readSizeLine:
		h.err = h.readSizeLine()
	case readIV:
		h.err = h.readIV()
	case decrypt:
		h.err = h.decrypt()
	case readInternalPayloadSize:
		h.err = h.readInternalPayloadSize()
	case readInternalTotalSize:
		h.err = h.readInternalTotalSize()
	case readZippedPayload:
		h.err = h.readZippedPayload()
	case readSignature:
		h.err = h.readSignature()
	case verifyHMAC:
		h.err = h.verifyHMAC()
	case invalid:
		h.err = errors.New(invalidMessage)
	}
	return h.err
}

func (h *Header) readMagicVersionLine() error {
	line, err := h.reader.ReadString('\n')
	if err != nil {
		return err
	}
	// find the magic "Syndie.Message.1." string
	if !strings.HasPrefix(line, syndieMessage) {
		return errors.New(invalidMessage)
	}
	return nil
}

func (h *Header) readHeaderKeyPairs() error {
	var counter int
	for {
		if counter > limit {
			return errors.New(invalidMessage)
		}
		line, err := h.reader.ReadString('\n')
		line = strings.TrimSpace(line)
		if err != nil {
			return err
		}
		if len(line) == 0 {
			break
		}
		h.ReadLine(line)
		counter++
	}
	return nil
}

func (h *Header) readSizeLine() error {
	line, err := h.reader.ReadString('\n')
	if err != nil {
		return errors.New(invalidMessage)
	}
	s, err := value(line)
	if err != nil {
		return errors.New(invalidMessage)
	}
	size, err := strconv.Atoi(s)
	if err != nil {
		return err
	}
	h.totalPayloadSize = size
	return nil
}

func (h *Header) decrypt() error {
	decrypted := make([]byte, h.totalPayloadSize)
	payload := make([]byte, h.totalPayloadSize)
	io.ReadFull(h.reader, payload)
	privKey, err := base64.I2PEncoding.DecodeString(h.BodyKey)
	if err != nil {
		return errors.New("error decoding: " + err.Error())
	}
	if h.totalPayloadSize%aes.BlockSize != 0 {
		return errors.New("ciphertext is not a multiple of the block size")
	}
	block, err := aes.NewCipher([]byte(privKey))
	if err != nil {
		return errors.New("error initializing NewCipher: %s" + err.Error())
	}
	decrypter := cipher.NewCBCDecrypter(block, h.iv)
	decrypter.CryptBlocks(decrypted, payload)
	h.encryptedPayload = payload
	h.enclosedReader = bytes.NewReader(decrypted)
	return nil
}

func (h *Header) readInternalPayloadSize() error {
	var counter int
	zero := make([]byte, 1)
	for {
		if counter > limit {
			return errors.New(invalidMessage)
		}
		_, err := h.enclosedReader.Read(zero)
		if err != nil {
			return errors.New(invalidMessage)
		}
		if bytes.Equal(zero, []byte{0x0}) {
			h.internalPayloadSize = readInt(h.enclosedReader)
			return nil
		}
		counter++
	}
}

func (h *Header) readInternalTotalSize() error {
	internalTotalSize := readInt(h.enclosedReader)
	if h.totalPayloadSize != internalTotalSize+len(h.iv) {
		return errors.New(invalidMessage)
	}
	return nil
}

func (h *Header) readIV() error {
	iv, err := h.reader.Peek(16)
	h.iv = iv
	return err
}

func (h *Header) readZippedPayload() error {
	var payload = make([]byte, h.internalPayloadSize)
	io.ReadFull(h.enclosedReader, payload)
	zr, err := zip.NewReader(bytes.NewReader(payload), int64(h.internalPayloadSize))
	if err != nil {
		return err
	}
	m, err := h.ParseMessage(zr)
	if err != nil {
		return err
	}
	h.msg = &m
	return nil
}

func (h *Header) readSignature() error {
	var err error
	h.signature, err = io.ReadAll(h.reader)
	return err
}

func (h *Header) verifyHMAC() error {
	scanner := bufio.NewScanner(bytes.NewBuffer(h.signature))
	scanner.Scan()
	// TODO: check authorization
	authorizationSig, err := value(scanner.Text())
	_ = authorizationSig
	if err != nil {
		return errors.New("invalid signature")
	}
	scanner.Scan()
	// TODO: check authentication
	authenticationSig, err := value(scanner.Text())
	if err != nil {
		return errors.New("invalid signature")
	}
	_ = authenticationSig
	// check the hmac
	var hmacPreKey bytes.Buffer
	k, _ := base64.I2PEncoding.DecodeString(h.BodyKey)
	hmacPreKey.Write(k)
	hmacPreKey.Write(h.iv)
	sha := sha256.New()
	sha.Write(hmacPreKey.Bytes())
	hm := hmac.New(sha256.New, sha.Sum(nil))
	hm.Write(h.encryptedPayload[16 : len(h.encryptedPayload)-32])
	if !hmac.Equal(hm.Sum(nil), h.encryptedPayload[len(h.encryptedPayload)-32:]) {
		return errors.New("unable to verify HMAC")
	}
	return nil
}
