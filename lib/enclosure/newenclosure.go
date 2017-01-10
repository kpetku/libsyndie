package enclosure

import (
	"bufio"
	"bytes"
	"errors"
	"os"
	"strconv"
	"strings"

	"github.com/kpetku/go-syndie/lib/syndieuri"

	log "github.com/Sirupsen/logrus"
)

func NewEnclosure(f *os.File, err error) (*SyndieHeader, *SyndieTrailer) {
	scanner := bufio.NewScanner(f)

	scanner.Split(bufio.SplitFunc(doubleNewlineDelimiter))

	success := scanner.Scan()
	encl := SyndieHeader{}

	if validateHeader(&encl, scanner.Bytes()) != nil {
		log.Infof("Error from validateHeader! %s", err)
	}
	scanner.Scan()
	tail := SyndieTrailer{}
	if validateBody(&tail, scanner.Bytes()) != nil {
		log.Infof("Error from validateBody! %s", err)
	}

	if serr := scanner.Err(); err != nil {
		log.Fatal(serr)
	}
	if !success {
		// False on error or EOF. Check error
		err = scanner.Err()
		if err == nil {
			log.Printf("OooF EOF")
		} else {
			log.Fatal(err)
		}
	}
	return &encl, &tail
}

func validateBody(e *SyndieTrailer, b []byte) error {
	if !bytes.Contains(b, []byte("Size=")) {
		return errors.New("Invalid Syndie body")
	}
	position := bytes.Index(b, []byte("\n"))
	if position > 0 {
		size, err := strconv.Atoi(string(b[5:position])) // 5 is len("Size=")

		if err != nil {
			log.WithFields(log.Fields{
				"at":     "(Enclosure) MarshallTrailer",
				"size":   size,
				"line":   string(b[5:position]),
				"reason": "parsing error",
			}).Fatalf("%s", err)
		}
		e.Size = size
		e.Raw = b[position+1:]
		e.Body = b

		e.AuthenticationSig = bytes.Split(b[bytes.Index(b, []byte("AuthenticationSig=")):], []byte("\n"))[0][18:]
		e.AuthorizationSig = bytes.Split(b[bytes.Index(b, []byte("AuthorizationSig=")):], []byte("\n"))[0][17:]
	} else {
		return errors.New("Flattened corrupted Syndie body...")
	}
	return nil
}

func validateHeader(e *SyndieHeader, b []byte) error {
	for count, line := range bytes.SplitN(b, []byte("\n"), 2) {
		if count == 0 {
			if !bytes.Contains(line, []byte("Syndie.Message.1.0")) {
				return errors.New("Invalid Syndie message")
			}
		} else {
			err := validateHeaderLine(e, line)
			if err != nil {
				log.Fatalf("%s", err.Error())
			}
		}
	}
	return nil
}

func validateHeaderLine(e *SyndieHeader, b []byte) error {
	if bytes.Contains(b, []byte("=")) {
		split := bytes.SplitN(b, []byte("="), 2)
		key := string(split[0])
		value := string(split[1])
		switch key {
		case "Author":
			e.Author = value
		case "AuthenticationMask":
			e.AuthenticationMask = value
		case "TargetChannel":
			e.TargetChannel = value
		case "PostURI":
			u := syndieuri.URI{}
			u.Marshall(value)
			e.PostURI = u
		case "References":
			var out []syndieuri.URI
			r := strings.Fields(value)
			for _, ref := range r {
				u := syndieuri.URI{}
				u.Marshall(ref)
				out = append(out, u)
			}
			e.References = out
		case "Tags":
			e.Tags = strings.Fields(value)
		case "OverwriteURI":
			u := syndieuri.URI{}
			u.Marshall(value)
			e.OverwriteURI = u
		case "ForceNewThread":
			if strings.Contains(value, "true") {
				e.ForceNewThread = true
			}
		case "RefuseReplies":
			if strings.Contains(value, "true") {
				e.RefuseReplies = true
			}
		case "Cancel":
			var out []syndieuri.URI
			r := strings.Fields(value)
			for _, canc := range r {
				u := syndieuri.URI{}
				u.Marshall(canc)
				out = append(out, u)
			}
			e.Cancel = out
		case "Subject":
			e.Subject = value
		case "BodyKey":
			e.BodyKey = value
		case "BodyKeyPromptSalt":
			e.BodyKeyPromptSalt = value
		case "BodyKeyPrompt":
			e.BodyKeyPrompt = value
		case "Identity":
			e.Identity = value
		case "EncryptKey":
			e.EncryptKey = value
		case "Name":
			e.Name = value
		case "Description":
			e.Description = value
		case "Edition":
			i, err := strconv.Atoi(value)
			if err != nil {
				log.WithFields(log.Fields{
					"at":     "(Enclosure) MarshallHeader strconv",
					"i":      i,
					"reason": "conversion error",
				}).Fatalf("%s", err)
			}
			e.Edition = i
		case "PublicPosting":
			if strings.Contains(value, "true") {
				e.PublicPosting = true
			}
		case "PublicReplies":
			if strings.Contains(value, "true") {
				e.PublicReplies = true
			}
		case "AuthorizedKeys":
			e.AuthorizedKeys = strings.Fields(value)
		case "ManagerKeys":
			e.ManagerKeys = strings.Fields(value)
		case "Archives":
			var out []syndieuri.URI
			r := strings.Fields(value)
			for _, arch := range r {
				u := syndieuri.URI{}
				u.Marshall(arch)
				out = append(out, u)
			}
			e.Archives = out
		case "ChannelReadKeys":
			e.ChannelReadKeys = strings.Fields(value)
		case "Expiration":
			e.Expiration = value
		case "Syndie.MessageType":
		// TODO: wrong place for MessageType?
		default:
			return errors.New("corrupt header key: " + key + " value: " + value)
		}
		return nil
	}
	return errors.New("corrupt header")
}

func doubleNewlineDelimiter(data []byte, atEOF bool) (advance int, token []byte, err error) {
	if atEOF && len(data) == 0 {
		return 0, nil, nil
	}
	if i := bytes.Index(data, []byte("\n\n")); i >= 0 {
		return i + 2, data[0:i], nil
	}
	if atEOF {
		return len(data), data, nil
	}
	return 0, nil, nil
}
