package enclosure

import (
	"bufio"
	"bytes"
	"os"
	"strconv"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/kpetku/go-syndie/lib/common"
	"github.com/kpetku/go-syndie/lib/syndieuri"
)

// Enclosure holds the reference to a Syndie Header and Message
type Enclosure struct {
	Header          *SyndieHeader
	Message         *SyndieTrailer
	isAuthenticated bool
	isAuthorized    bool
}

// OpenFile opens a file and returns a populated Enclosure
func (enclosure *Enclosure) OpenFile(s string) *Enclosure {
	var rest []byte
	buf := make([]byte, 0, 64*1024)

	file, err := os.Open(s)
	if err != nil {
		log.WithFields(log.Fields{
			"at":     "(main)",
			"file":   s,
			"reason": "failed to open file",
		}).Fatalf("%s", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Buffer(buf, 1024*1024)
	scanner.Split(common.NewlineDelimiter)

	for scanner.Scan() {
		bs := scanner.Bytes()
		if bytes.HasPrefix(bs, []byte("Syndie.Message.1.0")) {
			str := strings.SplitAfter(string(bs), "\n")
			enclosure.Header.Version = str[0]
			for _, h := range str {
				switch strings.Split(h, "=")[0] {
				case "Author":
					enclosure.Header.Author = strings.Split(h, "Author=")[1]
				case "AuthenticationMask":
					enclosure.Header.AuthenticationMask = strings.Split(h, "AuthenticationMask=")[1]
				case "TargetChannel":
					enclosure.Header.TargetChannel = strings.Split(h, "TargetChannel=")[1]
				case "PostURI":
					u := syndieuri.URI{}
					u.Marshall(strings.Split(h, "PostURI=")[1])
					enclosure.Header.PostURI = u
				case "References":
					var out []syndieuri.URI
					r := strings.Split(h, "References=")[1:]
					for _, ref := range r {
						u := syndieuri.URI{}
						u.Marshall(ref)
						out = append(out, u)
					}
					enclosure.Header.References = out
				case "Tags":
					enclosure.Header.Tags = strings.Split(h, "Tags=")[1:]
				case "OverwriteURI":
					u := syndieuri.URI{}
					u.Marshall(strings.Split(h, "OverwriteURI=")[1])
					enclosure.Header.OverwriteURI = u
				case "ForceNewThread":
					if strings.Contains(strings.Split(h, "ForceNewThread=")[1], "true") {
						enclosure.Header.ForceNewThread = true
					}
				case "RefuseReplies":
					if strings.Contains(strings.Split(h, "RefuseReplies=")[1], "true") {
						enclosure.Header.RefuseReplies = true
					}
				case "Cancel":
					var out []syndieuri.URI
					r := strings.Split(h, "Cancel=")[1:]
					for _, canc := range r {
						u := syndieuri.URI{}
						u.Marshall(canc)
						out = append(out, u)
					}
					enclosure.Header.Cancel = out
				case "Subject":
					enclosure.Header.Subject = strings.Split(h, "Subject=")[1]
				case "BodyKey":
					enclosure.Header.BodyKey = strings.Split(h, "BodyKey=")[1]
				case "BodyKeyPromptSalt":
					enclosure.Header.BodyKeyPromptSalt = strings.Split(h, "BodyKeyPromptSalt=")[1]
				case "BodyKeyPrompt":
					enclosure.Header.BodyKeyPrompt = strings.Split(h, "BodyKeyPrompt=")[1]
				case "Identity":
					enclosure.Header.Identity = strings.Split(h, "Identity=")[1]
				case "EncryptKey":
					enclosure.Header.EncryptKey = strings.Split(h, "EncryptKey=")[1]
				case "Name":
					enclosure.Header.Name = strings.Split(h, "Name=")[1]
				case "Description":
					enclosure.Header.Description = strings.Split(h, "Description=")[1]
				case "Edition":
					i, err := strconv.Atoi(strings.TrimRight(strings.Split(h, "=")[1], "\n"))
					if err != nil {
						log.WithFields(log.Fields{
							"at":     "(Enclosure) MarshallHeader strconv",
							"i":      i,
							"reason": "conversion error",
						}).Fatalf("%s", err)
					}
					enclosure.Header.Edition = i
				case "PublicPosting":
					if strings.Contains(strings.Split(h, "PublicPosting=")[1], "true") {
						enclosure.Header.PublicPosting = true
					}
				case "PublicReplies":
					if strings.Contains(strings.Split(h, "PublicReplies=")[1], "true") {
						enclosure.Header.PublicReplies = true
					}
				case "AuthorizedKeys":
					enclosure.Header.AuthorizedKeys = strings.Split(h, "AuthorizedKeys=")[1:]
				case "ManagerKeys":
					enclosure.Header.ManagerKeys = strings.Split(h, "ManagerKeys=")[1:]
				case "Archives":
					var out []syndieuri.URI
					r := strings.Split(h, "Archives=")[1:]
					for _, arch := range r {
						u := syndieuri.URI{}
						u.Marshall(arch)
						out = append(out, u)
					}
					enclosure.Header.Archives = out
				case "ChannelReadKeys":
					enclosure.Header.ChannelReadKeys = strings.Split(h, "ChannelReadKeys=")[1:]
				case "Expiration":
					enclosure.Header.Expiration = strings.Split(h, "Expiration=")[1]
				}
			}
		} else {
			rest = append(rest[:], bs...)
		}
		// TODO: err out here?
	}
	if err := scanner.Err(); err != nil {
		log.WithFields(log.Fields{
			"at":               "(Enclosure) scanner.Scan",
			"enclosure_parsed": scanner.Bytes(),
			"reason":           "invalid input scanned",
		}).Fatalf("%s", err)
	}
	if bytes.HasPrefix(rest, []byte("Size=")) {
		line := strings.Split(string(rest)[len("Size="):], "\n")
		size, err := strconv.Atoi(line[0])
		rest := strings.Join(line[1:], "")
		if err != nil {
			log.WithFields(log.Fields{
				"at":     "(Enclosure) MarshallTrailer",
				"size":   size,
				"line":   line,
				"reason": "parsing error",
			}).Fatalf("%s", err)
		}
		enclosure.Message.size = size
		enclosure.Message.raw = []byte(rest)
	} else {
		panic("Invalid trailer marshalling attempted")
	}
	return enclosure
}
