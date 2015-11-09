package handler

import (
	"dumbdumb"
	"errors"
	"fmt"
	"github.com/jmoiron/jsonq"
	"log"
	"napping"
	"regexp"
)

type TranslateHandler struct {
	GoogleAPIKey string
}

func (h TranslateHandler) BuildAPIParamsFromRequest(request dumbdumb.Request) (napping.Params, error) {
	// TODO store precompiled on handler or global
	exp := regexp.MustCompile(
		// e.g. "Translate en -> fr lobster"
		"^[Tt]ranslate\\s(?P<source>[a-z]{2})\\s->\\s(?P<target>[a-z]{2})\\s(?P<q>.*)$")
	match := exp.FindStringSubmatch(request.GetPayload())
	if match == nil || len(match) != 4 {
		return nil, errors.New("invalid request format")
	}
	params := napping.Params{
		"key":    h.GoogleAPIKey,
		"source": match[1],
		"target": match[2],
		"q":      match[3],
	}
	return params, nil
}

// TODO: test sending erroneous keys
// TODO: test 403 forbidden error handling (when setting up new Google API)
func (h TranslateHandler) HandleRequest(request dumbdumb.Request) error {
	log.Printf("TranslateHandler got request: %v", request.GetPayload())

	params, err := h.BuildAPIParamsFromRequest(request)
	if err != nil {
		log.Printf("TranslateHandler could not parse request '%v': %v",
			request.GetPayload(), err)
		return err
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

	// assume the first translation is the one we want
	//sourceLang, err := jq.String("data", "translations", "0", "detectedSourceLanguage")
	result, err := jq.String("data", "translations", "0", "translatedText")

	if err != nil {
		return err
	}

	err = request.SendOutput(fmt.Sprintf("%s -> %v: %v\n",
		params["source"], params["target"], result))

	return err
}
