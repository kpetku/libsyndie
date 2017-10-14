package syndieutil

import (
	"archive/zip"
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"strings"
)

type Message struct {
	Page       []Page
	Attachment []Attatchment
	Avatar     []byte
	References string
}

type Attatchment struct {
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
		key := string(split[0])
		value := string(split[1])
		switch key {
		case "Content-Type":
			p.ContentType = value
		case "Title":
			p.Title = value
		case "References":
			p.References = value
		default:
			return errors.New("malformed page")
		}
	}
	return nil
}

func (a *Attatchment) ReadLine(s string) error {
	if strings.Contains(s, "=") {
		split := strings.SplitN(s, "=", 2)
		key := string(split[0])
		value := string(split[1])
		switch key {
		case "Name":
			a.Name = value
		case "Content-Type":
			a.ContentType = value
		case "Description":
			a.Description = value
		default:
			return errors.New("malformed attachment")
		}
	}
	return nil
}

// ParseMessage is probably broken with multiple attachments/pages
func (h *Header) ParseMessage(zr *zip.Reader) (Message, error) {
	var pagenum, attachnum int
	m := Message{}
	for _, file := range zr.File {
		fileReader, err := file.Open()
		buf := new(bytes.Buffer)
		buf.ReadFrom(fileReader)

		newStr := buf.String()
		if err != nil {
			return Message{}, fmt.Errorf("Error opening enclosed zip file %s", err)

		}
		defer fileReader.Close()
		switch file.Name {
		case "headers.dat":
			scanner := bufio.NewScanner(strings.NewReader(newStr))
			for scanner.Scan() {
				h.ReadLine(scanner.Text())
			}
		case "references.cfg":
			m.References = newStr
		case "avatar32.png":
			m.Avatar = []byte(newStr)
		}
		if strings.HasPrefix(file.Name, "page") {
			p := Page{}
			m.Page = append(m.Page, p)
			if strings.HasSuffix(file.Name, ".dat") {
				m.Page[pagenum].Data = newStr
			}
			if strings.HasSuffix(file.Name, ".cfg") {
				scanner := bufio.NewScanner(strings.NewReader(newStr))
				for scanner.Scan() {
					m.Page[pagenum].ReadLine(scanner.Text())
				}
			}
			pagenum++
		}
		// The spec is unclear if this should be "attach" or "attachment"
		if strings.HasPrefix(file.Name, "attach") {
			a := Attatchment{}
			m.Attachment = append(m.Attachment, a)
			if strings.HasSuffix(file.Name, ".dat") {
				m.Attachment[attachnum].Data = []byte(newStr)
			}
			if strings.HasSuffix(file.Name, ".cfg") {
				scanner := bufio.NewScanner(strings.NewReader(newStr))
				for scanner.Scan() {
					m.Attachment[attachnum].ReadLine(scanner.Text())
				}
			}
			attachnum++
		}
	}
	return m, nil
}
