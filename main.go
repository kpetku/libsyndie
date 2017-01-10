package main

import (
	"os"

	log "github.com/Sirupsen/logrus"
)

func main() {
	log.Printf("go-syndie: startup.")
	_, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
}
