package handler

import (
	"errors"
	"fmt"
	"github.com/davemt/dumbdumb"
	"github.com/jmcvetta/napping"
	"github.com/jmoiron/jsonq"
	"log"
	"net/url"
	"regexp"
)

type TranslateHandler struct {
	GoogleAPIKey string
}

func (h TranslateHandler) BuildAPIParamsFromRequest(request dumbdumb.Request) (url.Values, error) {
	// TODO store precompiled on handler or global
	exp := regexp.MustCompile(
		// e.g. "Translate en -> fr lobster"
		"^[Tt]ranslate\\s(?P<source>[a-z]{2})\\s->\\s(?P<target>[a-z]{2})\\s(?P<q>.*)$")
	match := exp.FindStringSubmatch(request.GetPayload())
	if match == nil || len(match) != 4 {
		return nil, errors.New("invalid request format")
	}
	params := url.Values{}
	params.Add("key", h.GoogleAPIKey)
	params.Add("source", match[1])
	params.Add("target", match[2])
	params.Add("q", match[3])
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
