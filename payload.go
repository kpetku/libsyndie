package main

import (
	"bufio"
	"bytes"
	"os"
	"strconv"
	"strings"

	log "github.com/Sirupsen/logrus"
)

// SyndiePayload holds a Syndie Payload that contains a header and trailer field
type SyndiePayload struct {
	Header  []byte
	Trailer []byte
}

// Message holds the reference to a Syndie header and trailer
type Message struct {
	Head *SyndieHeader
	Body *SyndieTrailer
}

// MarshallTrailer marshalls a SyndiePayload into a SyndieTrailer with various struct fields
func (payload *SyndiePayload) MarshallTrailer() *SyndieTrailer {
	bs := payload.Trailer
	trailer := SyndieTrailer{}
	if bytes.HasPrefix(bs, []byte("Size=")) {
		line := strings.Split(string(bs)[len("Size="):], "\n")
		size, err := strconv.Atoi(line[0])
		rest := strings.Join(line[1:], "")
		if err != nil {
			log.WithFields(log.Fields{
				"at":     "(payload) MarshallTrailer",
				"size":   size,
				"line":   line,
				"reason": "parsing error",
			}).Fatalf("%s", err)
		}
		trailer.size = size
		trailer.body = []byte(rest)
	} else {
		panic("Invalid trailer marshalling attempted")
	}
	return &trailer
}

// MarshallHeader marshalls a SyndiePayload into a SyndieHeader with various struct fields
func (payload *SyndiePayload) MarshallHeader() *SyndieHeader {
	header := SyndieHeader{}
	str := strings.SplitAfter(string(payload.Header), "\n")
	header.Version = str[0]
	for _, h := range str {
		switch strings.Split(h, "=")[0] {
		case "Author":
			header.Author = strings.Split(h, "=")[1]
		case "AuthenticationMask":
			header.AuthenticationMask = strings.Split(h, "=")[1]
		case "TargetChannel":
			header.TargetChannel = strings.Split(h, "=")[1]
		case "PostURI":
			//			header.PostURI = strings.Split(h, "=")[1:]
		case "References":
			//			header.References = strings.Split(h, "=")[1:]
		case "Tags":
			header.Tags = strings.Split(h, "=")[1:]
		case "OverwriteURI":
			//			header.OverwriteURI = strings.Split(h, "=")[1:]
		case "ForceNewThread":
			if strings.Contains(strings.Split(h, "=")[1], "true") {
				header.ForceNewThread = true
			}
		case "RefuseReplies":
			if strings.Contains(strings.Split(h, "=")[1], "true") {
				header.RefuseReplies = true
			}
		case "Cancel":
			//			header.Cancel = strings.Split(h, "=")[1:]
		case "Subject":
			header.Subject = strings.Split(h, "=")[1]
		case "BodyKey":
			header.BodyKey = strings.Split(h, "=")[1]
		case "BodyKeyPromptSalt":
			header.BodyKeyPromptSalt = strings.Split(h, "=")[1]
		case "BodyKeyPrompt":
			header.BodyKeyPrompt = strings.Split(h, "=")[1]
		case "Identity":
			header.Identity = strings.Split(h, "=")[1]
		case "EncryptKey":
			header.EncryptKey = strings.Split(h, "=")[1]
		case "Name":
			header.Name = strings.Split(h, "=")[1]
		case "Description":
			header.Description = strings.Split(h, "=")[1]
		case "Edition":
			i, err := strconv.Atoi(strings.TrimRight(strings.Split(h, "=")[1], "\n"))
			if err != nil {
				log.WithFields(log.Fields{
					"at":     "(payload) MarshallHeader strconv",
					"i":      i,
					"reason": "conversion error",
				}).Fatalf("%s", err)
			}
			header.Edition = i
		case "PublicPosting":
			if strings.Contains(strings.Split(h, "=")[1], "true") {
				header.PublicPosting = true
			}
		case "PublicReplies":
			if strings.Contains(strings.Split(h, "=")[1], "true") {
				header.PublicReplies = true
			}
		case "AuthorizedKeys":
			header.AuthorizedKeys = strings.Split(h, "=")[1:]
		case "ManagerKeys":
			header.ManagerKeys = strings.Split(h, "=")[1:]
		case "Archives":
			//			header.Archives = strings.Split(h, "=")[1:]
		case "ChannelReadKeys":
			header.ChannelReadKeys = strings.Split(h, "=")[1:]
		case "Expiration":
			header.Expiration = strings.Split(h, "=")[1]
		}
	}
	return &header
}

// Parse marshalls a raw []byte representation of a Syndie Message into header and trailer fields
func (payload *SyndiePayload) Parse(bs []byte) {
	if bytes.HasPrefix(bs, []byte("Syndie.Message.1.0")) {
		payload.Header = bs
	} else {
		rest := append(payload.Trailer[:], bs...)
		payload.Trailer = rest
	}
}

//OpenFile tries to open a SyndiePayload
func (payload *SyndiePayload) OpenFile(s string) {
	file, err := os.Open(s)
	if err != nil {
		log.WithFields(log.Fields{
			"at":     "(main)",
			"file":   s,
			"reason": "failed to open file",
		}).Fatalf("%s", err)
	}
	defer file.Close()
	message := Message{}

	scanner := bufio.NewScanner(file)
	buf := make([]byte, 0, 64*1024)
	scanner.Buffer(buf, 1024*1024)
	scanner.Split(newlineDelimiter)

	for scanner.Scan() {
		payload.Parse(scanner.Bytes())
	}
	if err := scanner.Err(); err != nil {
		log.WithFields(log.Fields{
			"at":             "(main) scanner.Scan",
			"payload_parsed": scanner.Bytes(),
			"reason":         "invalid input scanned",
		}).Fatalf("%s", err)
	}

	message.Head = payload.MarshallHeader()
	message.Body = payload.MarshallTrailer()

	// TODO: Lookup the channel key from the URI and attempt to decrypt
	log.Printf("Dumping contents: %s", message.Body.DecryptAES(""))
}
