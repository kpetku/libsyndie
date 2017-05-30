package syndieutil

import (
	"bytes"
	"errors"
	"strconv"
	"strings"

	log "github.com/Sirupsen/logrus"
)

// MessageHeader holds a Syndie message header that contains version and pairs fields:
/*
Author AuthenticationMask TargetChannel PostURI References Tags OverwriteURI ForceNewThread
RefuseReplies Cancel Subject BodyKey BodyKeyPromptSalt BodyKeyPrompt Identity EncryptKey Name
Description Edition PublicPosting PublicReplies AuthorizedKeys ManagerKeys Archives ChannelReadKeys Expiration
*/
type MessageHeader struct {
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
	Archives           []URI
	ChannelReadKeys    []string
	Expiration         string
}

// TODO: clean up this ugliness
func validateHeaderLine(e *MessageHeader, b []byte) error {
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
			u := URI{}
			u.Marshall(value)
			e.PostURI = u
		case "References":
			var out []URI
			r := strings.Fields(value)
			for _, ref := range r {
				u := URI{}
				u.Marshall(ref)
				out = append(out, u)
			}
			e.References = out
		case "Tags":
			e.Tags = strings.Fields(value)
		case "OverwriteURI":
			u := URI{}
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
			var out []URI
			r := strings.Fields(value)
			for _, canc := range r {
				u := URI{}
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
			var out []URI
			r := strings.Fields(value)
			for _, arch := range r {
				u := URI{}
				u.Marshall(arch)
				out = append(out, u)
			}
			e.Archives = out
		case "ChannelReadKeys":
			e.ChannelReadKeys = strings.Fields(value)
		case "Expiration":
			e.Expiration = value
		// TODO: wrong place for MessageType?
		case "Syndie.MessageType":
		default:
			return errors.New("corrupt header key: " + key + " value: " + value)
		}
		return nil
	}
	return errors.New("corrupt header")
}
