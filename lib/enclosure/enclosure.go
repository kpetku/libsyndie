package enclosure

import (
	"bufio"
	"bytes"
	"os"
	"strconv"
	"strings"

	"github.com/kpetku/go-syndie/lib/common"
	"github.com/kpetku/go-syndie/lib/syndieuri"

	log "github.com/Sirupsen/logrus"
)

// Enclosure holds the reference to a Syndie Header and Message
type Enclosure struct {
	Header            *SyndieHeader
	Message           *SyndieTrailer
	AuthenticationSig string
	AuthorizationSig  string
	HmacPos           int
}

// OpenFile opens a file and returns a populated Enclosure
func (enclosure *Enclosure) OpenFile(s string) *Enclosure {
	var rest2 []byte

	file, err := os.Open(s)
	if err != nil {
		log.WithFields(log.Fields{
			"at":     "(main)",
			"file":   s,
			"reason": "failed to open file",
		}).Fatalf("%s", err)
	}
	defer file.Close()
	stat, err := file.Stat()
	if err != nil {
		log.WithFields(log.Fields{
			"at":     "(main)",
			"file":   s,
			"reason": "failed stat file",
		}).Fatalf("%s", err)
	}
	length := int(stat.Size())

	log.Printf("File size is %d", length)
	buf := make([]byte, 0, length)
	scanner := bufio.NewScanner(file)
	scanner.Buffer(buf, length)
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
			if bytes.Index(bs, []byte("AuthorizationSig=")) > 0 {
				apos := bytes.Index(bs, []byte("AuthorizationSig="))
				sig := strings.Split(string(bs[apos:]), "AuthorizationSig=")
				enclosure.HmacPos = apos
				enclosure.AuthorizationSig = strings.TrimSpace(sig[1])
			}
			if bytes.Index(bs, []byte("AuthenticationSig=")) > 0 {
				apos := bytes.Index(bs, []byte("AuthenticationSig="))
				sig := strings.Split(string(bs[apos:]), "AuthenticationSig=")
				enclosure.AuthenticationSig = strings.TrimSpace(sig[1])
			}
			rest2 = append(rest2[:], bs...)
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
	if bytes.HasPrefix(rest2, []byte("Size=")) {
		line := strings.Split(string(rest2)[len("Size="):], "\n")
		size, err := strconv.Atoi(line[0])
		rest := strings.Join(line[1:], "\n")
		if err != nil {
			log.WithFields(log.Fields{
				"at":     "(Enclosure) MarshallTrailer",
				"size":   size,
				"line":   line,
				"reason": "parsing error",
			}).Fatalf("%s", err)
		}
		enclosure.Message.Size = size
		enclosure.Message.Raw = []byte(rest)
	} else {
		panic("Invalid trailer marshalling attempted")
	}
	return enclosure
}
