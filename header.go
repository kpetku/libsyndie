package main

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
