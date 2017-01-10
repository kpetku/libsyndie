package common

import (
	"archive/zip"
	"bytes"

	log "github.com/Sirupsen/logrus"
)

func ParseZipReader(content []byte, size int) {
	zr, err := zip.NewReader(bytes.NewReader(content), int64(size))
	if err != nil {
		panic("error")
	}
	var page, attatchment int
	for _, file := range zr.File {
		fileReader, err := file.Open()
		buf := new(bytes.Buffer)
		buf.ReadFrom(fileReader)
		newStr := buf.String()
		if err != nil {
			log.Fatalf("Error opening enclosed zip file %s", err)
		}
		defer fileReader.Close()
		switch file.Name[:4] {
		case "atta":
			log.Printf("Attatchment num %d: %s", attatchment, newStr)
			attatchment++
		case "page":
			log.Printf("Page num %d: %s", page, newStr)
			page++
		case "avat":
			log.Printf("avatar32.png: %s", newStr)
		case "refe":
			log.Printf("references.cfg: %s", newStr)
		case "head":
			log.Printf("Headers.dat: %s", newStr)
		}
	}
}
