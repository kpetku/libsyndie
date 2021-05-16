package client

import (
	"encoding/binary"
	"errors"
	"io"
	"strconv"

	"github.com/go-i2p/go-i2p/lib/common/base64"
	"github.com/kpetku/libsyndie/archive"
)

const upperBoundLimit = 10000
const invalidArchiveServer = "invalid syndie archive server"

type Client struct {
	ChannelHashes []archive.ChannelHash
	Messages      []archive.Message
	Urls          []string
	archive.Header
}

type reader struct {
	r   io.Reader
	err error
}

func (r *reader) read(data interface{}) {
	if r.err == nil {
		r.err = binary.Read(r.r, binary.BigEndian, data)
	}
}

func New() *Client {
	return &Client{}
}

func (c *Client) Parse(input io.Reader) error {
	var url []string
	r := reader{r: input}

	// Read ArchiveFlags (unimplemented)
	r.read(&c.ArchiveFlags)

	// Read the admin channel
	r.read(&c.AdminChannel)

	// Count the number of alternate URIs
	r.read(&c.NumAltURIs)
	if int(c.NumAltURIs) > upperBoundLimit {
		return errors.New(invalidArchiveServer + ": too many alternate archive URIs")
	}

	// Populate AltURIs with other known archive servers
	var archiveAltURIs []string
	for i := 0; i < int(c.NumAltURIs); i++ {
		var length uint16
		r.read(&length)

		uri := make([]byte, int(length))
		r.read(&uri)

		archiveAltURIs = append(archiveAltURIs, string(uri))
	}
	c.AltURIs = archiveAltURIs

	// Count the number of channels
	r.read(&c.NumChannels)
	if int(c.NumChannels) > upperBoundLimit {
		return errors.New(invalidArchiveServer + ": too many channels")
	}

	// Read the channel hashes
	for i := 0; i < int(c.NumChannels); i++ {
		var hash archive.ChannelHash
		r.read(&hash)
		url = append(url, base64.I2PEncoding.EncodeToString(hash.ChannelHash[:])+"/meta.syndie")
		c.ChannelHashes = append(c.ChannelHashes, hash)
	}

	r.read(&c.NumMessages)
	if int(c.NumMessages) > upperBoundLimit {
		return errors.New(invalidArchiveServer + ": too many messages")
	}

	// Read messages and append urls
	var message archive.Message
	for i := 0; i < int(c.NumMessages); i++ {
		r.read(&message)
		url = append(url, base64.I2PEncoding.EncodeToString(c.ChannelHashes[int(message.ScopeChannel)].ChannelHash[:])+"/"+strconv.Itoa(int(message.MessageID))+".syndie")
		c.Messages = append(c.Messages, message)
	}
	if r.err != nil {
		return r.err
	}
	c.Urls = url
	return nil
}
