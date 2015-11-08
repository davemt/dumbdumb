package handler

import (
	"dumbdumb"
	"fmt"
	"github.com/jmoiron/jsonq"
	"log"
	"napping"
	"strings"
)

type TranslateHandler struct {
	GoogleAPIKey string
}

func (h TranslateHandler) HandleRequest(request dumbdumb.Request) error {
	log.Printf("TranslateHandler got request: %v", request.GetPayload())
	parts := strings.SplitAfterN(request.GetPayload(), " ", 2)
	_, query := parts[0], parts[1]

	params := napping.Params{
		"key": h.GoogleAPIKey,
		"q":   query,
		// TODO: add support for to/from languages
		"language": "en",
	}

	var data map[string]interface{}

	// make request to get place id based on text query
	resp, err := napping.Get(
		"https://www.googleapis.com/language/translate/v2",
		&params, &data, nil)
	if resp.Status() != 200 {
		return err
	}

	jq := jsonq.NewQuery(data)

	// assume the first match is the best
	result, err := jq.String("data", "translations", "0", "translatedText")

	if err != nil {
		return err
	}

	err = request.SendOutput(fmt.Sprintf("English: %v\n",
		result))

	return err
}
