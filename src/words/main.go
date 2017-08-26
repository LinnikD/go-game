package words

import (
	"time"
	"net/http"
	"log"
)

func CheckWordExists(word string) {
	for {
		_, err := http.Get("http://example.com/")
		if err != nil {
			log.Printf("get words error: %s", err)
		}
		time.Sleep(time.Minute)
	}
}


