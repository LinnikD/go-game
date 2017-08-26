package words

import (
	"time"
	"net/http"
	"log"
)

type Word struct {
	ID       int64    `bson:"channel_id,omitempty" json:"channel_id,omitempty"`
	Text     string   `bson:"channel,omitempty" json:"channel,omitempty"`

}

func CheckWordExists(word string) {
	for {

		// check mongo word
		// if not ok check yandex
		// write yandex result
		// return

		_, err := http.Get("http://example.com/")
		if err != nil {
			log.Printf("get words error: %s", err)
		}
		time.Sleep(time.Minute)
	}
}


