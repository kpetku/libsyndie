package syndieutil

import (
	"bytes"
	"errors"
	"fmt"
	"strings"

	"github.com/jackpal/bencode-go"
)

/*
URI defines the URIs safely passable within syndie, capable of referencing specific resources.
They contain one of four reference types, plus a bencoded set of attributes:
*/
type URI struct {
	RefType     string   `bencode:"-"`
	Name        string   `bencode:"name,omitempty"`
	Desc        string   `bencode:"desc,omitempty"`
	Tag         []string `bencode:"tag,omitempty"`
	Author      string   `bencode:"author,omitempty"`
	Net         string   `bencode:"net,omitempty"`
	ReadKeyType string   `bencode:"readKeyType,omitempty"`
	ReadKeyData string   `bencode:"readKeyData,omitempty"`
	PostKeyType string   `bencode:"postKeyType,omitempty"`
	PostKeyData string   `bencode:"postKeyData,omitempty"`
	URL         string   `bencode:"url,omitempty"`
	Channel     string   `bencode:"channel,omitempty"`
	MessageID   int      `bencode:"messageId,omitempty"`
	Page        int      `bencode:"page,omitempty"`
	Attachment  int      `bencode:"attatchment,omitempty"`
	Scope       []string `bencode:"scope,omitempty"`
	PostByScope []string `bencode:"postbyscope,omitempty"`
	Age         int      `bencode:"age,omitempty"`
	AgeLocal    int      `bencode:"agelocal,omitempty"`
	UnreadOnly  bool     `bencode:"unreadonly,omitempty"`
	TagInclude  []string `bencode:"taginclude,omitempty"`
	TagRequire  []string `bencode:"tagrequire,omitempty"`
	TagExclude  []string `bencode:"tagexclude,omitempty"`
	TagMessages bool     `bencode:"tagmessages,omitempty"`
	PageMin     int      `bencode:"pagemin,omitempty"`
	PageMax     int      `bencode:"pagemax,omitempty"`
	AttachMin   int      `bencode:"attachmin,omitempty"`
	AttachMax   int      `bencode:"attachmax,omitempty"`
	RefMin      int      `bencode:"refmin,omitempty"`
	RefMax      int      `bencode:"refmax,omitempty"`
	KeyMin      int      `bencode:"keymin,omitempty"`
	KeyMax      int      `bencode:"keymax,omitempty"`
	Encrypted   bool     `bencode:"encrypted,omitempty"`
	PBE         bool     `bencode:"pbe,omitempty"`
	Private     bool     `bencode:"private,omitempty"`
	Public      bool     `bencode:"public,omitempty"`
	Authorized  bool     `bencode:"authorized,omitempty"`
	Threaded    bool     `bencode:"threaded,omitempty"`
	Keyword     string   `bencode:"keyword,omitempty"`
	Body        string   `bencode:"body,omitempty"`
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
func (u *URI) Marshall(s string) error {
	if len(s) < 3 {
		return errors.New("URI was too short to process")
	}
	u.RefType = strings.Split(s, ":")[0]
	prepared, err := prepareURI(s)
	if err != nil {
		return err
	}
	r := bytes.NewReader([]byte(prepared))

	berr := bencode.Unmarshal(r, &u)
	if berr != nil {
		return fmt.Errorf("error while parsing bencode: %s", berr)
	}
	return nil
}
func (u *URI) String() string {
	var buf []byte
	w := bytes.NewBuffer(buf)
	bencode.Marshal(w, *u)
	out, _ := prepareURI(w.String())
	return "urn:syndie:" + u.RefType + ":" + out
}
