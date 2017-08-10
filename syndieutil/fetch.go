package syndieutil

import (
	"io/ioutil"
	"log"
)

func FetchFromDisk(path string) {
	var totalChan, totalMsg, totalMeta int
	fetchChannelList, _ := ioutil.ReadDir(path)
	for _, c := range fetchChannelList {
		totalChan++
		log.Printf("Channel found: %s\n", c.Name())
		fetchMessageList, _ := ioutil.ReadDir(path + c.Name())
		for _, m := range fetchMessageList {
			if m.Name() == "meta.syndie" {
				totalMeta++
			} else {
				totalMsg++
			}
		}
	}
	log.Printf("Fetched: %d channels, %d messages, %d meta.", totalChan, totalMsg, totalMeta)
}
