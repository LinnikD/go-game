package words

import (
	"net/http"
	"log"
	"model"
	"fmt"
)

type WordChecker struct {
	api      string
	mongo    *model.Connection
}

func NewWordChecker(api string, mongo *model.Connection) (*WordChecker) {
	checker := WordChecker{api, mongo}
	return &checker
}

func (c *WordChecker) CheckWordExists(text string) (bool) {

	// check mongo word
	word := c.mongo.GetWord(text)

	if word != nil {
		return word.IsWord
	}

	// if not ok check yandex
	wordCheckUrl := fmt.Sprintf(c.api, text)
	response, err := http.Get(wordCheckUrl)
	if err != nil {
		log.Printf("get words error: %s", err)
	}

	print(response)
	// write yandex result

	// return
	return true
}


