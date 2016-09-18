package enclosure

import (
	"github.com/kpetku/go-syndie/lib/syndieuri"
)

// SyndieHeader holds a Syndie header that contains version and pairs fields:
/*
Author AuthenticationMask TargetChannel PostURI References Tags OverwriteURI ForceNewThread
RefuseReplies Cancel Subject BodyKey BodyKeyPromptSalt BodyKeyPrompt Identity EncryptKey Name
Description Edition PublicPosting PublicReplies AuthorizedKeys ManagerKeys Archives ChannelReadKeys Expiration
*/
type SyndieHeader struct {
	Version            string
	Author             string
	AuthenticationMask string
	TargetChannel      string
	PostURI            syndieuri.URI
	References         []syndieuri.URI
	Tags               []string
	OverwriteURI       syndieuri.URI
	ForceNewThread     bool
	RefuseReplies      bool
	Cancel             []syndieuri.URI
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
	Archives           []syndieuri.URI
	ChannelReadKeys    []string
	Expiration         string
}
