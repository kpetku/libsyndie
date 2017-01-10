package syndieuri

import (
	"bytes"
	"errors"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/jackpal/bencode-go"
)

/*
URI defines the URIs safely passable within syndie, capable of referencing specific resources.
They contain one of four reference types, plus a bencoded set of attributes:
*/
type URI struct {
	RefType     string
	Name        string   `bencode:",omitempty"`
	Desc        string   `bencode:",omitempty"`
	Tag         []string `bencode:",omitempty"`
	Author      string   `bencode:",omitempty"`
	Net         string   `bencode:",omitempty"`
	ReadKeyType string   `bencode:",omitempty"`
	ReadKeyData string   `bencode:",omitempty"`
	PostKeyType string   `bencode:",omitempty"`
	PostKeyData string   `bencode:",omitempty"`
	URL         string   `bencode:",omitempty"`
	Channel     string   `bencode:",omitempty"`
	MessageID   int      `bencode:",omitempty"`
	Page        int      `bencode:",omitempty"`
	Attachment  int      `bencode:",omitempty"`
	Scope       []string `bencode:",omitempty"`
	PostByScope []string `bencode:",omitempty"`
	Age         int      `bencode:",omitempty"`
	AgeLocal    int      `bencode:",omitempty"`
	UnreadOnly  bool     `bencode:",omitempty"`
	TagInclude  []string `bencode:",omitempty"`
	TagRequire  []string `bencode:",omitempty"`
	TagExclude  []string `bencode:",omitempty"`
	TagMessages bool     `bencode:",omitempty"`
	PageMin     int      `bencode:",omitempty"`
	PageMax     int      `bencode:",omitempty"`
	AttachMin   int      `bencode:",omitempty"`
	AttachMax   int      `bencode:",omitempty"`
	RefMin      int      `bencode:",omitempty"`
	RefMax      int      `bencode:",omitempty"`
	KeyMin      int      `bencode:",omitempty"`
	KeyMax      int      `bencode:",omitempty"`
	Encrypted   bool     `bencode:",omitempty"`
	PBE         bool     `bencode:",omitempty"`
	Private     bool     `bencode:",omitempty"`
	Public      bool     `bencode:",omitempty"`
	Authorized  bool     `bencode:",omitempty"`
	Threaded    bool     `bencode:",omitempty"`
	Keyword     string   `bencode:",omitempty"`
	Body        string   `bencode:",omitempty"`
}

// colonDelimiter returns data for a scanner delimited by a colon
func colonDelimiter(data []byte, atEOF bool) (advance int, token []byte, err error) {
	log.Printf("data was %s", data)
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

// trimSyndieURI trims "urn:", "urn:syndie", and "syndie:" from the left of a string
func trimSyndieURI(in string) string {
	if strings.HasPrefix(in, "urn:syndie:") {
		in = strings.Join(strings.Split(in, "urn:syndie:")[1:], "")
	}
	if strings.HasPrefix(in, "urn:") {
		in = strings.Join(strings.Split(in, "urn:")[1:], "")
	}
	if strings.HasPrefix(in, "syndie:") {
		in = strings.Join(strings.Split(in, "syndie:")[1:], "")
	}
	return in
}

// prepareURI checks a string for an invalid URI and ommits the refType before bencoding the rest.
func prepareURI(in string) (out string, err error) {
	if len(in) < 3 {
		return in, errors.New("invalid URI")
	}
	in = trimSyndieURI(in)
	switch strings.Split(strings.ToLower(in), ":")[0] {
	case "url", "channel", "search", "archive", "text":
		// Drop the RefType to prepare for bencode
		return strings.Join(strings.Split(in, ":")[1:], ":"), nil
	default:
		return in, errors.New("invalid URI refType: " + in)
	}
}

// Marshall takes a URI as string and returns a populated URI
func (u *URI) Marshall(s string) *URI {
	if len(s) < 3 {
		log.WithFields(log.Fields{
			"at":     "(uri) Marshall",
			"reason": "URI was too short to process",
		}).Fatalf("URI was too short to process")
		return &URI{}
	}
	s = trimSyndieURI(s)
	u.RefType = strings.Split(s, ":")[0]
	prepared, err := prepareURI(s)
	if err != nil {
		panic(err)
	}
	r := bytes.NewReader([]byte(prepared))

	berr := bencode.Unmarshal(r, &u)
	if berr != nil {
		log.WithFields(log.Fields{
			"at":     "(uri) Marshall",
			"reason": "error while parsing bencode",
		}).Infof("%s", berr)
		panic(err)
	}
	return u
}
