package main

import (
	"bytes"
	"encoding/base64"
)

// TODO: depend on go-i2p/common

const alphabet = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789-~"

// I2PEncoding returns an I2P compatible base64 encoding based on a custom alphabet
var I2PEncoding *base64.Encoding = base64.NewEncoding(alphabet)

// newlineDelimiter returns data for a scanner delimited by \n\n
func newlineDelimiter(data []byte, atEOF bool) (advance int, token []byte, err error) {
	if atEOF && len(data) == 0 {
		return 0, nil, nil
	}
	if i := bytes.Index(data, []byte{'\n', '\n'}); i >= 0 {
		return i + 2, data[0:i], nil
	}
	if atEOF {
		return len(data), data, nil
	}
	return 0, nil, nil
}
