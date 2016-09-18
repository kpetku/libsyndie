package main

import (
	"os"

	log "github.com/Sirupsen/logrus"
)

func main() {
	log.Printf("go-syndie: startup.")
	payload := SyndiePayload{}
	payload.OpenFile(os.Args[1])
}
