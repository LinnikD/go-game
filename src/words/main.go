package words

import (
	"net/http"
	"log"
	"model"
	"fmt"
	"io/ioutil"
	"encoding/json"
	"net/url"
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
	parameters := url.Values{}
	parameters.Add("text", text)
	p := parameters.Encode()

	wordCheckUrl := fmt.Sprintf(c.api, p)
	response := apiRequest(wordCheckUrl)

	// write yandex result
	if len(response) > 0 {
		c.mongo.UpsertWord(model.Word{text, false})
		return false
	} else {
		c.mongo.UpsertWord(model.Word{text, true})
		return true
	}
}

func apiRequest(url string) ([]apiResponse) {
	response, err := http.Get(url)

	if err != nil {
		log.Printf("%s", err)
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Printf("%s", err)
	}

	respBody := make([]apiResponse,0)
	err = json.Unmarshal(body, &respBody)

	return respBody
}

type apiResponse struct {
	Code int `json:"code,omitempty"`
}