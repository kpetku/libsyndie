package main

import (
	"flag"
	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/kpetku/go-syndie/syndieutil"
)

func main() {
	var arg string
	if len(os.Args) > 1 {
		arg = os.Args[1]
	}

	defaultKey := flag.String("key", "pjvUqwqXVD5Da7pJPVJcYStnfBrWaPqQCPN8Jw8Q-Lw=", "use specified key, default: pjvUqwqXVD5Da7pJPVJcYStnfBrWaPqQCPN8Jw8Q-Lw=")

	flag.Parse()

	file, err := os.Open(arg)
	if err != nil {
		log.Fatalf("Error while opening file %s", err)
	}

	header := syndieutil.New(syndieutil.BodyKey(*defaultKey))
	message, err := header.Unmarshal(file)
	if err != nil {
		log.Printf("Error reading message: %s", err.Error())
	}
	log.Printf("Opened message: %s", message.Page[0].Data)
	defer file.Close()
}
