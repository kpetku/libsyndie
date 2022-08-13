package archive

type Archive struct {
	ChannelHashes []ChannelHash
	Messages      []Message
	Urls          []string
	Header
}

type Header struct {
	ArchiveFlags uint16
	AdminChannel uint32
	AltURIs      []string
	NumAltURIs   byte
	NumChannels  uint32
	NumMessages  uint32
}

type ChannelHash struct {
	ChannelHash    [32]byte
	ChannelEdition uint64
	ChannelFlags   byte
}

type Message struct {
	MessageID     uint64
	ScopeChannel  uint32
	TargetChannel uint32
	MsgFlags      byte
}
