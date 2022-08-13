package syndieutil

import (
	"strconv"
	"strings"

	"github.com/kpetku/libsyndie/crypto"
)

const metaMessageType string = "Syndie.MessageType=meta"

type Metadata struct {
	Identity   *crypto.SigningKeypair
	EncryptKey *crypto.PrivateReplyKeypair
	BodyKey    crypto.SessionKey
	Edition    int
	Name       string
}

func NewMetadata() *Metadata {
	return new(Metadata)
}

// New creates a new metadata "channel" with new signing/reply key pairs
func (m *Metadata) New(name string) error {
	m.Name = name
	// Create the "BodyKey" a temporary session key for the encrypted message itself
	sess := buildSessionKey()
	m.BodyKey = sess
	// Create the "Identity" aka manager/manager-pub (signing) key pairs
	skp, err := buildSigningKeypair()
	if err != nil {
		return err
	}
	m.Identity = skp
	// Create the "EncryptKey" aka reply/reply-pub (private reply) key pairs
	e, err := buildPrivateReplyKeypair()
	if err != nil {
		return err
	}
	m.EncryptKey = e
	// Create the first edition
	m.Edition = buildEdition()
	return nil
}

func (m Metadata) String() string {
	var h Header
	h.Set(Name(m.Name))
	h.Set(BodyKey(m.BodyKey))
	h.Set(Identity(m.Identity.String()))
	h.Set(Edition(m.Edition))
	h.Set(EncryptKey(m.EncryptKey.String()))
	h.Set(BodyKey(crypto.NewSessionKey()))

	var sb strings.Builder
	sb.WriteString(syndieMessage)
	// TODO: Version this properly
	sb.WriteString("0")
	sb.WriteString(newLine)
	if h.Name != "" {
		sb.WriteString("Name=")
		sb.WriteString(m.Name)
		sb.WriteString(newLine)
	}
	if h.BodyKey != "" {
		sb.WriteString("BodyKey=")
		sb.WriteString(h.BodyKey)
		sb.WriteString(newLine)
	}
	if h.Edition != 0 {
		sb.WriteString("Edition=")
		sb.WriteString(strconv.Itoa(h.Edition))
		sb.WriteString(newLine)
	}
	if h.EncryptKey != "" {
		sb.WriteString("EncryptKey=")
		sb.WriteString(h.EncryptKey)
		sb.WriteString(newLine)
	}
	if h.Identity != "" {
		sb.WriteString("Identity=")
		sb.WriteString(h.Identity)
		sb.WriteString(newLine)
	}
	sb.WriteString(metaMessageType)
	sb.WriteString(newLine)
	return sb.String()
}

func buildSessionKey() string {
	return crypto.NewSessionKey()
}

func buildSigningKeypair() (*crypto.SigningKeypair, error) {
	skp := crypto.NewSigningKeypair()
	err := skp.Generate()
	if err != nil {
		return &crypto.SigningKeypair{}, err
	}
	return skp, nil
}

func buildPrivateReplyKeypair() (*crypto.PrivateReplyKeypair, error) {
	r, err := crypto.NewPrivateReplyKeypair()
	if err != nil {
		return &crypto.PrivateReplyKeypair{}, err
	}
	return r, nil
}

func buildEdition() int {
	// TODO: actually generate this
	return 1
}
