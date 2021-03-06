package handler

import (
	"fmt"
	"github.com/davemt/dumbdumb"
	"github.com/jmcvetta/napping"
	"github.com/jmoiron/jsonq"
	"log"
	"net/url"
	"strings"
)

type PlaceDirectoryHandler struct {
	GoogleAPIKey string
}

func (h PlaceDirectoryHandler) HandleRequest(request dumbdumb.Request) error {
	log.Printf("PlaceDirectoryHandler got request: %v", request.GetPayload())
	parts := strings.SplitAfterN(request.GetPayload(), " ", 2)
	_, placeQuery := parts[0], parts[1]

	params := url.Values{}
	params.Add("key", h.GoogleAPIKey)
	params.Add("query", placeQuery)

	var data map[string]interface{}

	// make request to get place id based on text query
	resp, err := napping.Get(
		"https://maps.googleapis.com/maps/api/place/textsearch/json",
		&params, &data, nil)
	if resp.Status() != 200 {
		return err
	}

	jq := jsonq.NewQuery(data)

	// assume the first match is the best
	placeId, err := jq.String("results", "0", "place_id")

	// use place id to query place details
	params = url.Values{}
	params.Add("key", h.GoogleAPIKey)
	params.Add("placeid", placeId)

	resp, err = napping.Get(
		"https://maps.googleapis.com/maps/api/place/details/json",
		&params, &data, nil)
	if resp.Status() != 200 {
		return err
	}

	jq = jsonq.NewQuery(data)

	result, err := jq.Object("result")
	if err != nil {
		return err
	}
	placeName := result["name"].(string)
	phoneNumber, ok := result["international_phone_number"]
	phoneNumberStr := "unknown"
	if ok {
		phoneNumberStr = phoneNumber.(string)
	}
	address := result["formatted_address"].(string)
	isOpen, err := jq.Bool("result", "opening_hours", "open_now")
	openNowStr := "unknown"
	if err == nil {
		if isOpen {
			openNowStr = "yes"
		} else {
			openNowStr = "no"
		}
	}
	err = request.SendOutput(fmt.Sprintf("%v\nPhone: %v\nAddr: %v\nOpen now: %v",
		placeName, phoneNumberStr, address, openNowStr))
	return err
}
