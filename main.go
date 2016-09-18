package main

import (
	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/kpetku/go-syndie/lib/enclosure"
)

func main() {
	log.Printf("go-syndie: startup.")
	e := enclosure.Enclosure{Header: &enclosure.SyndieHeader{}, Message: &enclosure.SyndieTrailer{}}
	e.OpenFile(os.Args[1])
}
