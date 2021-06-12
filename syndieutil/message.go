package syndieutil

import (
	"archive/zip"
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"strings"
)

type Message struct {
	Page       []Page
	Attachment []Attachment
	Avatar     []byte
	References string
}

type Attachment struct {
	Name        string
	ContentType string
	Description string
	Data        []byte
}

type Page struct {
	ContentType string
	Title       string
	References  string
	Data        string
}

func (p *Page) ReadLine(s string) error {
	if strings.Contains(s, "=") {
		split := strings.SplitN(s, "=", 2)
		key := strings.ToLower(string(split[0]))
		value := string(split[1])
		switch key {
		case "content-type":
			p.ContentType = value
		case "title":
			p.Title = value
		case "references":
			p.References = value
		default:
			return errors.New("malformed page")
		}
	}
	return nil
}

// ParseMessage is probably broken with multiple attachments/pages
func (h *Header) ParseMessage(zr *zip.Reader) (Message, error) {
	var foundPageCfg, foundPageDat bool
	var foundAttachCfg, foundAttachDat bool
	var pagenum, attachnum int
	attachnum++
	m := Message{}
	a := Attachment{}
	for _, file := range zr.File {
		fileReader, err := file.Open()
		if err != nil {
			return Message{}, fmt.Errorf("error opening enclosed zip file %s", err)
		}
		contents, err := ioutil.ReadAll(fileReader)
		if err != nil {
			return Message{}, fmt.Errorf("error reading from enclosed zip file %s", err)
		}
		defer fileReader.Close()
		switch file.Name {
		case "headers.dat":
			scanner := bufio.NewScanner(bytes.NewReader(contents))
			for scanner.Scan() {
				h.ReadLine(scanner.Text())
			}
		case "references.cfg":
			m.References = string(contents)
		case "avatar32.png":
			m.Avatar = contents
		}
		if strings.HasPrefix(file.Name, "page") {
			if foundPageCfg && foundPageDat {
				pagenum++
			}
			p := Page{}
			m.Page = append(m.Page, p)
			if strings.HasSuffix(file.Name, ".dat") {
				foundPageDat = true
				m.Page[pagenum].Data = string(contents)
			}
			if strings.HasSuffix(file.Name, ".cfg") {
				foundPageCfg = true
				scanner := bufio.NewScanner(bytes.NewReader(contents))
				for scanner.Scan() {
					m.Page[pagenum].ReadLine(scanner.Text())
				}
			}
		}
		// The spec is unclear if this should be "attach" or "attachment"
		if strings.HasPrefix(file.Name, "attachment") {
			if strings.HasSuffix(file.Name, ".dat") {
				a.Data = contents
				foundAttachDat = true
				continue
			}
			if strings.HasSuffix(file.Name, ".cfg") {
				scanner := bufio.NewScanner(bytes.NewReader(contents))
				for scanner.Scan() {
					if strings.Contains(scanner.Text(), "=") {
						split := strings.SplitN(scanner.Text(), "=", 2)
						key := strings.ToLower(string(split[0]))
						value := string(split[1])
						switch key {
						case "name":
							a.Name = value
						case "content-type":
							a.ContentType = value
						case "description":
							a.Description = value
						}
					}
				}
				foundAttachCfg = true
				if foundAttachCfg && foundAttachDat {
					m.Attachment = append(m.Attachment, a)
					foundAttachCfg = false
					foundAttachDat = false
				}
				continue
			}
		}
	}
	return m, nil
}
