package main

import (
	"flag"
	"os"

	log "github.com/Sirupsen/logrus"

	"github.com/kpetku/go-syndie/syndieutil"
)

func main() {
	var arg string
	if len(os.Args) > 0 {
		arg = os.Args[1]
	}

	defaultKey := flag.String("key", "pjvUqwqXVD5Da7pJPVJcYStnfBrWaPqQCPN8Jw8Q-Lw=", "use specified key, default: pjvUqwqXVD5Da7pJPVJcYStnfBrWaPqQCPN8Jw8Q-Lw=")

	flag.Parse()

	file, err := os.Open(arg)
	if err != nil {
		log.Fatalf("Error while opening file %s", err)
	}

	err2 := syndieutil.ParseBody(file, *defaultKey)

	if err2 != nil {
		log.Printf("Error reading message: %s", err2.Error())
	}

	defer file.Close()

	/*
		current, err := user.Current()
		if err != nil {
			log.Fatalf("Could not obtain current user")
		}
		syndieutil.FetchFromDisk(current.HomeDir + "/.syndie/archive/")
	*/

}
