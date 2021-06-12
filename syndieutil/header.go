package syndieutil

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"strconv"
	"strings"
)

// Header holds a Syndie message header that contains version and pairs fields
type Header struct {
	Version            string
	Author             string
	AuthenticationMask string
	TargetChannel      string
	PostURI            URI
	References         []URI
	Tags               []string
	OverwriteURI       URI
	ForceNewThread     bool
	RefuseReplies      bool
	Cancel             []URI
	Subject            string
	BodyKey            string
	BodyKeyPromptSalt  string
	BodyKeyPrompt      string
	Identity           string
	EncryptKey         string
	Name               string
	Description        string
	Edition            int
	PublicPosting      bool
	PublicReplies      bool
	AuthorizedKeys     []string
	ManagerKeys        []string
	Archives           string
	ChannelReadKeys    string
	Expiration         string
	MessageType        string

	reader              *bufio.Reader
	state               state
	err                 error
	iv                  []byte
	enclosedReader      *bytes.Reader
	encryptedPayload    []byte
	totalPayloadSize    int
	internalPayloadSize int
	msg                 *Message
	signature           []byte
}

type state int

const (
	readMagicVersionLine state = iota
	readHeaderKeyPairs
	readSizeLine
	readIV
	decrypt
	readInternalPayloadSize
	readInternalTotalSize
	readZippedPayload
	readSignature
	verifyHMAC
	invalid
)

// New creates a new Header and accepts a list of option functions
func New(opts ...func(*Header)) *Header {
	h := &Header{}

	// call option functions on instance to set options on it
	for _, opt := range opts {
		opt(h)
	}

	return h
}

// Set sets the specified option functions
func (h *Header) Set(opts ...func(*Header)) *Header {
	// call option functions on instance to set options on it
	for _, opt := range opts {
		opt(h)
	}

	return h
}

// ReadLine takes a key=value pair and reads it into the current header
func (h *Header) ReadLine(s string) error {
	if strings.Contains(s, "=") {
		split := strings.SplitN(s, "=", 2)
		key := string(split[0])
		value := strings.Join(split[1:], " ")
		switch key {
		case "Author":
			h.Set(Author(value))
		case "AuthenticationMask":
			h.Set(AuthenticationMask(value))
		case "TargetChannel":
			h.Set(TargetChannel(value))
		case "PostURI":
			h.Set(PostURI(parseSingleURI(value)))
		case "References":
			h.Set(References(parseSliceURI(value)))
		case "Tags":
			h.Set(Tags(parseSliceString(value)))
		case "OverwriteURI":
			h.Set(OverwriteURI(parseSingleURI(value)))
		case "ForceNewThread":
			h.Set(ForceNewThread(parseBool(value)))
		case "RefuseReplies":
			h.Set(RefuseReplies(parseBool(value)))
		case "Cancel":
			h.Set(Cancel(parseSliceURI(value)))
		case "Subject":
			h.Set(Subject(value))
		case "BodyKey":
			h.Set(BodyKey(value))
		case "BodyKeyPromptSalt":
			h.Set(BodyKeyPromptSalt(value))
		case "BodyKeyPrompt":
			h.Set(BodyKeyPrompt(value))
		case "Identity":
			h.Set(Identity(value))
		case "EncryptKey":
			h.Set(EncryptKey(value))
		case "Name":
			h.Set(Name(value))
		case "Description":
			h.Set(Description(value))
		case "Edition":
			i, err := strconv.Atoi(value)
			if err != nil {
				return fmt.Errorf("conversion error: %s", err)
			}
			h.Set(Edition(i))
		case "PublicPosting":
			h.Set(PublicPosting(parseBool(value)))
		case "PublicReplies":
			h.Set(PublicReplies(parseBool(value)))
		case "AuthorizedKeys":
			h.Set(AuthorizedKeys(parseSliceString(value)))
		case "ManagerKeys":
			h.Set(ManagerKeys(parseSliceString(value)))
		case "Archives":
			h.Set(Archives(value))
		case "ChannelReadKeys":
			h.Set(ChannelReadKeys(value))
		case "Expiration":
			h.Set(Expiration(value))
		case "Syndie.MessageType":
			h.Set(MessageType(value))
		default:
			return errors.New("unknown header")
		}
		return nil
	}
	return errors.New("malformed header")
}

// Author is an optional function of Header
func Author(author string) func(*Header) {
	return func(h *Header) {
		h.Author = author
	}
}

// AuthenticationMask is an optional function of Header
func AuthenticationMask(authenticationmask string) func(*Header) {
	return func(h *Header) {
		h.AuthenticationMask = authenticationmask
	}
}

// TargetChannel is an optional function of Header
func TargetChannel(targetchannel string) func(*Header) {
	return func(h *Header) {
		h.TargetChannel = targetchannel
	}
}

// PostURI is an optional function of Header
func PostURI(postURI URI) func(*Header) {
	return func(h *Header) {
		h.PostURI = postURI
	}
}

// References is an optional function of Header
func References(references []URI) func(*Header) {
	return func(h *Header) {
		h.References = references
	}
}

// Tags is an optional function of Header
func Tags(tags []string) func(*Header) {
	return func(h *Header) {
		h.Tags = tags
	}
}

// OverwriteURI is an optional function of Header
func OverwriteURI(overwriteURI URI) func(*Header) {
	return func(h *Header) {
		h.OverwriteURI = overwriteURI
	}
}

// ForceNewThread is an optional function of Header
func ForceNewThread(forcenewthread bool) func(*Header) {
	return func(h *Header) {
		h.ForceNewThread = forcenewthread
	}
}

// RefuseReplies is an optional function of Header
func RefuseReplies(refusereplies bool) func(*Header) {
	return func(h *Header) {
		h.RefuseReplies = refusereplies
	}
}

// Cancel is an optional function of Header
func Cancel(cancel []URI) func(*Header) {
	return func(h *Header) {
		h.Cancel = cancel
	}
}

// Subject is an optional function of Header
func Subject(subject string) func(*Header) {
	return func(h *Header) {
		h.Subject = subject
	}
}

// BodyKey is an optional function of Header
func BodyKey(bodykey string) func(*Header) {
	return func(h *Header) {
		h.BodyKey = bodykey
	}
}

// BodyKeyPromptSalt is an optional function of Header
func BodyKeyPromptSalt(bodykeypromptsalt string) func(*Header) {
	return func(h *Header) {
		h.BodyKeyPromptSalt = bodykeypromptsalt
	}
}

// BodyKeyPrompt is an optional function of Header
func BodyKeyPrompt(bodykeyprompt string) func(*Header) {
	return func(h *Header) {
		h.BodyKeyPrompt = bodykeyprompt
	}
}

// Identity is an optional function of Header
func Identity(identity string) func(*Header) {
	return func(h *Header) {
		h.Identity = identity
	}
}

// EncryptKey is an optional function of Header
func EncryptKey(encryptkey string) func(*Header) {
	return func(h *Header) {
		h.EncryptKey = encryptkey
	}
}

// Name is an optional function of Header
func Name(name string) func(*Header) {
	return func(h *Header) {
		h.Name = name
	}
}

// Description is an optional function of Header
func Description(description string) func(*Header) {
	return func(h *Header) {
		h.Description = description
	}
}

// Edition is an optional function of Header
func Edition(edition int) func(*Header) {
	return func(h *Header) {
		h.Edition = edition
	}
}

// PublicPosting is an optional function of Header
func PublicPosting(publicposting bool) func(*Header) {
	return func(h *Header) {
		h.PublicPosting = publicposting
	}
}

// PublicReplies is an optional function of Header
func PublicReplies(publicreplies bool) func(*Header) {
	return func(h *Header) {
		h.PublicReplies = publicreplies
	}
}

// AuthorizedKeys is an optional function of Header
func AuthorizedKeys(authorizedkeys []string) func(*Header) {
	return func(h *Header) {
		h.AuthorizedKeys = authorizedkeys
	}
}

// ManagerKeys is an optional function of Header
func ManagerKeys(managerkeys []string) func(*Header) {
	return func(h *Header) {
		h.ManagerKeys = managerkeys
	}
}

// Archives is an optional function of Header
func Archives(archives string) func(*Header) {
	return func(h *Header) {
		h.Archives = archives
	}
}

// ChannelReadKeys is an optional function of Header
func ChannelReadKeys(channelreadkeys string) func(*Header) {
	return func(h *Header) {
		h.ChannelReadKeys = channelreadkeys
	}
}

// Expiration is an optional function of Header
func Expiration(expiration string) func(*Header) {
	return func(h *Header) {
		h.Expiration = expiration
	}
}

// MessageType is an optional function of Header
func MessageType(messagetype string) func(*Header) {
	return func(h *Header) {
		h.MessageType = messagetype
	}
}

func parseSliceURI(value string) []URI {
	var out []URI
	r := strings.Fields(value)
	for _, arch := range r {
		u := URI{}
		u.Marshall(arch)
		out = append(out, u)
	}
	return out
}

func parseSingleURI(value string) URI {
	out := URI{}
	out.Marshall(value)
	return out
}

func parseBool(value string) bool {
	return value == "true"
}
func parseSliceString(value string) []string {
	return strings.Fields(value)
}
